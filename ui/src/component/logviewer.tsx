import useLogs from "../hook/useLogs";

const LogViewer = () => {
    const {logs} = useLogs();

    return (
    <div className="w-1/4 bg-gray-50 text-gray-800 text-xs px-3 py-1 overflow-scroll">
        {logs.map((log) => <span>{log}</span>)}        
    </div>
  )
}

export default LogViewer
