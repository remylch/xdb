import { axiosInstance } from ".";
import {CollectionData, CollectionsResponse} from "../type/collection";

export const fetchCollections = async (): Promise<CollectionsResponse> =>
  (await axiosInstance.get("/collections")).data;

export const fetchCollectionData = async (collection: string): Promise<CollectionData> =>
    (await axiosInstance.get(`/collections/${collection}`)).data;

export const createCollection = async (name: string): Promise<void> =>
    (await axiosInstance.post(`/collections?name=${name}`)).data;