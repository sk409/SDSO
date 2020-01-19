FROM centos:7

RUN yum -y update \
    && yum -y install epel-release \
    && yum -y install http://rpms.remirepo.net/enterprise/remi-release-7.rpm \
    && yum -y upgrade \
    && yum -y install nginx \
    && curl -sL https://rpm.nodesource.com/setup_9.x | bash - \
    && yum install -y gcc-c++ make \
    && yum install -y nodejs \
    && npm install -g vue-cli \
    && npm install -g nuxt \
    && npm install -g create-nuxt-app \
    && mkdir /sdso

ENV HOST 0.0.0.0
WORKDIR /sdso
CMD ["nginx", "-g", "daemon off;"]