export type Config = {
  rpc_bind_addr: string;
  rpc_public_addr: string;
  s3_api?: S3API;
  s3_web?: S3Web;
  admin?: Admin;
};

export type Admin = {
  api_bind_addr: string;
};

export type S3API = {
  s3_region: string;
  api_bind_addr: string;
  root_domain: string;
};

export type S3Web = {
  bind_addr: string;
  root_domain: string;
  index: string;
};
