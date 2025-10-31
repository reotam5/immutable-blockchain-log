import Settings from '@/components/settings'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/settings')({
  component: RouteComponent,
})

function RouteComponent() {
  return <Settings />
}
