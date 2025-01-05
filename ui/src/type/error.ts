import {AxiosError} from "axios";

export type ErrorResponse = AxiosError<{error: string}>