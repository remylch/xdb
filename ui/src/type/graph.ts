export type NodeInfoResponse = {
  version: string;
  ip: string;
  localStorePath: string;
  localLogPath: string;
};

export type GraphInfoResponse = {
  nodes: string[];
  clients: string[];
};
