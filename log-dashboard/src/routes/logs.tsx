import Logs from '@/components/logs'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/logs')({
  component: RouteComponent,
})

function RouteComponent() {
  return <Logs />
}
