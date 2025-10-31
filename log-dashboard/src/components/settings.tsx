import { useState } from "react"
import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { Label } from "./ui/label"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { API_BASE_URL } from "@/constants"
import { Alert, AlertDescription, AlertTitle } from "./ui/alert"
import { AlertCircleIcon } from "lucide-react"

function Settings() {
  const queryClient = useQueryClient()

  const query = useQuery({
    queryKey: ['log-path'],
    queryFn: async () => {
      return await (await fetch(`${API_BASE_URL}/settings/log`)).json()
    }
  })

  const mutation = useMutation({
    mutationFn: async () => {
      const res = await fetch(`${API_BASE_URL}/settings/log`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ path: editingPath })
      })

      if (!res.ok) {
        const errorData = await res.json().catch(() => ({}))
        const message = errorData.error
        throw new Error(message)
      }

      return res.json()
    },
    onSuccess: () => {
      console.log('Log path updated successfully')
      queryClient.invalidateQueries({ queryKey: ['log-path'] })
      setIsEditing(false)
    },
    onError: (error) => {
      setErrorMessage(error.message)
    }
  })

  const [editingPath, setEditingPath] = useState<string | null>(null)
  const [isEditing, setIsEditing] = useState(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  const onToggleEdit = () => {
    if (isEditing) {
      setIsEditing(false)
    } else {
      setEditingPath(query.data?.path)
      setErrorMessage(null)
      setIsEditing(true)
    }
  }

  return (
    <div className="space-y-6 max-w-7xl">
      <div className="">
        <Label htmlFor="log-path" className="mb-2 block">
          Monitoring Log Path:
        </Label>
        {
          isEditing ? (
            <div className="space-y-4">
              <div className="flex space-x-2">
                <Input
                  value={editingPath ?? ""}
                  onChange={(e) => setEditingPath(e.target.value)}
                />
                <Button onClick={() => mutation.mutate()}>
                  Save
                </Button>
              </div>
              {
                errorMessage && (
                  <Alert variant="destructive">
                    <AlertCircleIcon />
                    <AlertTitle>Unable to save log-path</AlertTitle>
                    <AlertDescription>
                      <p>{errorMessage}</p>
                    </AlertDescription>
                  </Alert>
                )
              }
            </div>
          ) : (
            <div className="flex space-x-2 items-center">
              <Button onClick={onToggleEdit}>
                Edit
              </Button>
              <div>{query.data?.path}</div>
            </div>
          )
        }
      </div>
    </div >
  )
}

export default Settings
