import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { API_BASE_URL } from "@/constants";
import { useInfiniteQuery } from "@tanstack/react-query";
import { Button } from "./ui/button";

interface Log {
  Timestamp: string;
  Content: string;
  IsValid: boolean;
  Source: string;
}

interface LogsResponse {
  logs: Log[];
  bookmark: string | null;
  hasNextPage: boolean;
}


function Logs() {
  const fetchLogs = async ({ pageParam = "" }): Promise<LogsResponse> => {
    const res = await fetch(`${API_BASE_URL}/log?filter=gateway-client&pageSize=10&bookmark=` + pageParam)
    return res.json()
  }

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isFetchingNextPage,
  } = useInfiniteQuery({
    initialPageParam: "",
    queryKey: ['logs'],
    queryFn: fetchLogs,
    getNextPageParam: (lastPage) => lastPage.hasNextPage ? lastPage.bookmark : null,
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
          {isFetching && !data?.pages?.length ? (
            <div>
              Loading...
            </div>
          ) : (
            data?.pages?.length === 0 ? (
              <div>
                No logs found
              </div>
            ) : (
              data?.pages?.flatMap(page => page.logs)?.map((log, index) => (
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
      <div>
        <Button disabled={isFetchingNextPage || !hasNextPage} onClick={() => {
          fetchNextPage()
        }}>
          Load More
        </Button>
      </div>
    </div>
  );
}

export default Logs;
