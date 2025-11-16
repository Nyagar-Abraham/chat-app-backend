# Create a private network in the cloud
# with ip range 10.0.0.0 - 10.0.255.255
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = {
    Name        = "${var.app_name}-vpc"
    Environment = var.environment
  }
}


# Create an igw to allow the vpc to connect to the internet
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  
  tags = {
    Name = "${var.app_name}-igw"
  }
}

# Create a public subnet in the vpc
# in two azs for high availability
resource "aws_subnet" "public" {
  count                   = 2
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.${count.index + 1}.0/24"
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  
  tags = {
    Name = "${var.app_name}-public-${count.index + 1}"
  }
}

# Create a private subnet in the vpc
# for app service , database and backend services
resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 10}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]
  
  tags = {
    Name = "${var.app_name}-private-${count.index + 1}"
  }
}

# Static public IP addresses
# One for each NAT Gateway
# Persist even if NAT Gateway is recreated
resource "aws_eip" "nat" {
  count  = 2
  domain = "vpc"
  
  tags = {
    Name = "${var.app_name}-nat-eip-${count.index + 1}"
  }
}


# Allows private subnet resources to access the internet (for updates, API calls, etc.)
# Receives incoming responses but blocks inbound connections
# High availability: One per AZ, so if one AZ fails, the other continues working
# Placed in public subnets
resource "aws_nat_gateway" "main" {
  count         = 2
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
  
  tags = {
    Name = "${var.app_name}-nat-${count.index + 1}"
  }
}

# All internet traffic (0.0.0.0/0) goes to the Internet Gateway
# Shared by both public subnets
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
  
  tags = {
    Name = "${var.app_name}-public-rt"
  }
}


# All internet traffic goes through the NAT Gateway
# Separate route table per private subnet → each uses its own AZ's NAT Gateway
resource "aws_route_table" "private" {
  count  = 2
  vpc_id = aws_vpc.main.id
  
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main[count.index].id
  }
  
  tags = {
    Name = "${var.app_name}-private-rt-${count.index + 1}"
  }
}

resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Links subnets to their respective route tables
# Determines how traffic is routed from each subnet
resource "aws_route_table_association" "private" {
  count          = 2
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

# Automatically fetches available AZs in your region
# Makes the code portable across regions
data "aws_availability_zones" "available" {
  state = "available"
}
