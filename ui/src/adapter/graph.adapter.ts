import { axiosInstance } from ".";
import {GraphResponse} from "../type/graph.ts";

export const fetchGraphInfos = async (): Promise<GraphResponse> =>
    (await axiosInstance.get("/graph")).data;