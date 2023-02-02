# #build
# FROM golang AS builder
# WORKDIR /go/src/
# COPY go.* ./
# COPY aws-vmcreate.go .
# RUN go mod download
# RUN go build -o aws-vmcreate

#deploy
FROM registry.access.redhat.com/ubi8/ubi-minimal

# Labels
LABEL name="aws-vmcreate" \
    maintainer="pandurang.com" \
    vendor="pandurang" \
    version="1.0.0" \
    release="1" \
    summary="This service enables provisioning/de-provisioning of AWS cloud vms." \
    description="This service enables provisioning/de-provisioning AWS cloud vms."


COPY aws-vmcreate aws-vmcreate

ENV ec2_tag_key "POC"
ENV ec2_tag_value "PandurangGolangOperator"
ENV ec2_command "create"
ENV ec2_instance_type "t2.micro"
ENV ec2_image_id "ami-0d0ca2066b861631c"

CMD ["bash","-c","./aws-vmcreate -c  $ec2_command -n $ec2_tag_key -v $ec2_tag_value -i $ec2_image_id -t $ec2_instance_type "]

