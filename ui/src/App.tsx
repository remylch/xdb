import { useCollections } from './hook/useCollections.ts';
import React, { PropsWithChildren } from "react";
import { useNavigate } from "react-router-dom";
import {Loader} from "./component/loader.tsx";

function App() {
  const { location, collections, error, isLoading } = useCollections()
    const navigate = useNavigate()
  return (
    <div className="flex flex-col gap-5 h-screen flex-1">
        {error ? <span>Unable to load collections</span> : null}
        {isLoading && <Loader text="Loading collections" /> }
        {location ? <span>Node data directory : {location}</span> : null}
        {collections.length > 0 ?
            <div className="flex flex-wrap gap-5">
                {collections?.map((collection) =>
                    <CollectionCard key={collection} onClick={() => navigate(`/${collection}`)}>
                        {collection}
                    </CollectionCard>)}
            </div>

            : <span>No collections found</span>}
    </div>
  );
}

const CollectionCard: React.FC<PropsWithChildren<{onClick: () => void }>> = ({children, onClick}) =>
     <div onClick={onClick} className="px-5 py-2 border rounded-md w-fit flex items-center gap-5 cursor-pointer">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="w-5 h-5">
            <path strokeLinecap="round" strokeLinejoin="round"
                  d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z"/>
        </svg>
        {children}
    </div>



export default App;
