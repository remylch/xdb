import {useQuery} from "react-query";
import {GraphResponse} from "../type/graph.ts";
import {fetchGraphInfos} from "../adapter/graph.adapter.ts";

export function useGraph() {
    const {data, error, isFetching, isLoading} = useQuery<GraphResponse>(
        ["/graph"],
        {
            queryFn: fetchGraphInfos,
            staleTime: 1000 * 60 * 60 * 60,
        }
    );

    return {}
}