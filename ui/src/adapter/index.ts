import axios from "axios";

export const axiosInstance = axios.create({
  baseURL: `${import.meta.env.VITE_NODE_HOST}:${
    import.meta.env.VITE_NODE_PORT
  }`,
});
