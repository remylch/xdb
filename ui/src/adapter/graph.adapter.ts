import { axiosInstance } from ".";
import { GraphInfoResponse } from "../type/graph.ts";

export const fetchGraphInfos = async (): Promise<GraphInfoResponse> =>
  (await axiosInstance.get("/graph")).data;
