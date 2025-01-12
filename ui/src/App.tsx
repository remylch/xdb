import {useCollections} from './hook/useCollections.ts';
import {FC, FormEvent, PropsWithChildren} from "react";
import {Loader} from "./component/loader.tsx";
import {useNavigate} from "react-router-dom";
import {createCollection} from "./adapter/collections.adapter.ts";
import {useQueryClient} from "react-query";
import {Title} from "./component/title.tsx";
import {TextInput} from "./component/input.tsx";
import {Button} from "./component/button.tsx";
import {ErrorResponse} from "./type/error.ts";
import { useNode } from './hook/useNode.ts';
import { Error } from './component/error.tsx';
import { useGraph } from './hook/useGraph.ts';

function App() {
    return (
        <div className="flex flex-col gap-5 h-screen flex-1">
            <Infos/>
            <CollectionModule/>
        </div>
    );
}

const Infos: FC = () =>
    <div className="flex flex-wrap gap-5 w-full">
        <NodeInfos/>
        <GraphInfos/>
    </div>

const NodeInfos: FC = () => {
    const {data, isFetching, isLoading, error} = useNode()

    if (isLoading || isFetching) {
        return <Loader text='Loading node info...' />
    }

    if (error) {
        return <Error errorType='internal' text='Error retrieving node infos' />
    }

    return <div className="p-5 border border-gray-900 rounded-md flex-1">
        <Title>Node Info</Title>
        <ul>
            <li>Version: <span className="text-gray-600">{data?.version}</span></li>
            <li>IP: <span className="text-gray-600">{data?.ip}</span></li>
            <li>Storage path: <span className="text-gray-600">{data?.localStorePath}</span></li>
            <li>Log path: <span className="text-gray-600">{data?.localLogPath}</span></li>
        </ul>
    </div>
}

const GraphInfos: FC = () => {
    const {data, error, isFetching, isLoading} = useGraph()

    if (isLoading || isFetching) {
        return <Loader text='Loading graph info...' />
    }

    if (error) {
        return <Error errorType='internal' text='Error retrieving graph infos' />
    }

    return <div className="p-5 border border-gray-900 rounded-md flex-1">
        <Title>Graph Info</Title>
        <ul>
            <li>Active nodes: <span className="text-gray-600">{data?.nodes}</span></li>
            <li>Clients: <span className="text-gray-600">{data?.clients}</span></li>
        </ul>
    </div>
}

const CollectionModule = () => {
    return (
        <div className="flex flex-col gap-5">
            <Title>Collections</Title>
            <CreateCollectionForm/>
            <Collections/>
        </div>
    )
}

const Collections: FC = () => {
    const navigate = useNavigate()
    const {collections, error, isLoading} = useCollections()

    if (isLoading) {
        return <Loader text="Loading collections"/>
    }

    if (error) {
        return <span>An error appear while trying to get collections.</span>
    }

    if (collections.length === 0) {
        return <span>No collections found</span>
    }

    return <div className="grid lg:grid-cols-4 gap-5 md:grid-cols-3 grid-cols-1">
        {collections.map((collection) =>
            <CollectionCard key={collection.name} onClick={() => navigate(`/${collection}`)}>
                {collection.name}
            </CollectionCard>)}
    </div>
}

const CreateCollectionForm: FC = () => {
    const queryClient = useQueryClient()

    const handleCreateCollection = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        const form = e.currentTarget
        const formData = new FormData(form)
        const collectionName = formData.get('name') as string

        if (collectionName.trim() === "") {
            alert("Collection name cannot be empty")
            return
        }

        await createCollection(collectionName).then(() => {
            queryClient.invalidateQueries("/collections")
            alert("Collection created successfully")
            form.reset()
        }).catch((err: ErrorResponse) => alert(`An error occurred while trying to create the collection : ${err.response?.data.error}`))
    }

    return (
        <form className="flex gap-5" onSubmit={handleCreateCollection}>
            <TextInput type="text" name="name" placeholder="Collection name"/>
            <Button type="submit">Create</Button>
        </form>
    )
}

const CollectionCard: FC<PropsWithChildren<{ onClick: () => void }>> = ({children, onClick}) =>
    <div onClick={onClick} className="px-5 py-2 border rounded-md flex items-center gap-5 cursor-pointer">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor"
             className="w-5 h-5">
            <path strokeLinecap="round" strokeLinejoin="round"
                  d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z"/>
        </svg>
        {children}
    </div>


export default App;
