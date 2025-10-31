import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { API_BASE_URL } from "@/constants";
import { useQuery } from "@tanstack/react-query";

interface Log {
  Timestamp: string;
  Content: string;
  IsValid: boolean;
  Source: string;
}


function Logs() {
  const query = useQuery({
    queryKey: ['logs'],
    queryFn: async () => {
      return await (await fetch(`${API_BASE_URL}/log`)).json() as Log[]
    }
  })

  return (
    <div className="max-w-7xl">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Timestamp</TableHead>
            <TableHead>Content</TableHead>
            <TableHead>Verified</TableHead>
            <TableHead>Client</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {query.isLoading ? (
            <div>
              Loading...
            </div>
          ) : (
            query?.data?.length === 0 ? (
              <div>
                No logs found
              </div>
            ) : (
              query?.data?.map((log, index) => (
                <TableRow key={index}>
                  <TableCell>{new Date(log.Timestamp).toLocaleString()}</TableCell>
                  <TableCell>{log.Content}</TableCell>
                  <TableCell>{log.IsValid ? 'Yes' : 'No'}</TableCell>
                  <TableCell>{log.Source}</TableCell>
                </TableRow>
              ))
            )
          )}
        </TableBody>
      </Table>
    </div>
  );
}

export default Logs;
