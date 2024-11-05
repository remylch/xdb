import {useQuery} from "react-query";
import {CollectionData} from "../type/collection.ts";
import {fetchCollectionData} from "../adapter/collections.adapter.ts";

export const useCollectionData = (collection: string) => {
    const { data, error, isLoading, isFetching } = useQuery<CollectionData>([`/collections/${collection}`, collection], {
        queryFn: () => fetchCollectionData(collection),
        staleTime: 1000 * 60 * 60 * 60,
    })

    return {
        isLoading: isLoading || isFetching,
        error,
        data,
    }
}
