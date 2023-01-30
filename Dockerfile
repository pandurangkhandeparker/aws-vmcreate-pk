FROM registry.access.redhat.com/ubi8/go-toolset:1.18.4-8.1669838000

# Labels
LABEL name="aws-vmcreate" \
    maintainer="xyzcompany.com" \
    vendor="xyzcompany" \
    version="1.0.0" \
    release="1" \
    summary="This service enables provisioning/de-provisioning of AWS cloud vms." \
    description="This service enables provisioning/de-provisioning AWS cloud vms."

# copy code to the build path
USER root
WORKDIR /opt
RUN chgrp -R 0 /opt && \
    chmod -R g=u /opt && \
    chmod +x -R /opt
USER 1001

ENV ec2_tag_key "POC"
ENV ec2_tag_value "GolangOperator"
ENV ec2_command "create"
ENV ec2_instance_type "t2.micro"
ENV ec2_image_id "ami-0d0ca2066b861631c"
#can be delete also

# ARG tag
# ENV ec2_tag_key=$key

# ARG tag_val
# ENV ec2_tag_value=$val

# ARG cmd
# ENV ec2_command=$cmd

COPY go.* ./
COPY aws-vmcreate.go .
RUN go mod download
RUN go build -o aws-vmcreate
CMD ["bash","-c","/opt/aws-vmcreate -c  $ec2_command -n $ec2_tag_key -v $ec2_tag_value -i $ec2_image_id -t $ec2_instance_type "]
