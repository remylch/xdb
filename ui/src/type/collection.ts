export type CollectionsResponse = {
  collections: Collection[];
};

export type Collection = {
  id: string;
  name: string;
  indexes: string[];
};

export type Data = unknown;
