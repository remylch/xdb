import {useParams} from "react-router-dom";
import {useCollectionData} from "../hook/useCollectionData.ts";
import {Loader} from "../component/loader.tsx";

export const CollectionPage = () => {
    const {collection} = useParams()

    const {error, isLoading, data} = useCollectionData(collection!)

    return <div className="flex flex-col gap-5">
        <span>{collection}</span>
        {error ? <span>Unable to load collection data</span> : null}
        {isLoading && <Loader text={`Loading ${collection} collection data`}/>}
        {data && <div>{data.data}</div>}
    </div>
}