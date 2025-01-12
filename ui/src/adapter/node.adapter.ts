import { axiosInstance } from ".";
import { NodeInfoResponse } from "../type/graph.ts";

export const fetchNodeInfos = async (): Promise<NodeInfoResponse> =>
  (await axiosInstance.get("/node")).data;
