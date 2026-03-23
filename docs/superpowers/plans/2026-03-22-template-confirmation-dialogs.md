# Template Confirmation Dialogs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add confirmation dialogs to Duplicate and Archive buttons on the template detail page, with name customization for Duplicate.

**Architecture:** Add optional `name` field to the duplicate API endpoint (OpenAPI → Go handler), then add two confirmation dialogs on the frontend: a `Dialog` with name input for Duplicate, and an `AlertDialog` for Archive. Unarchive remains immediate (low-risk restore action).

**Tech Stack:** OpenAPI 3.0, Go (oapi-codegen), Next.js, Radix UI Dialog/AlertDialog, TanStack Query

---

### Task 1: Add optional request body to duplicate endpoint in OpenAPI spec

**Files:**
- Modify: `api/openapi/openapi.yaml:128-146`

- [ ] **Step 1: Add DuplicateProgramTemplateRequest schema and requestBody**

In `api/openapi/openapi.yaml`, add a request body to the duplicate endpoint and a new schema:

Under `/program-templates/{id}/duplicate` post, add `requestBody` (between `tags` and `responses`):

```yaml
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/DuplicateProgramTemplateRequest'
```

Under `components.schemas`, add:

```yaml
    DuplicateProgramTemplateRequest:
      type: object
      properties:
        name:
          type: string
          maxLength: 200
          description: Custom name for the duplicate. If omitted, defaults to "{original name} (copy)".
```

Also update the endpoint description:

```yaml
      description: Creates a deep copy of the ProgramTemplate with all entries. Optionally accepts a custom name; defaults to appending " (copy)".
```

- [ ] **Step 2: Regenerate Go code**

Run from `api/`:
```bash
task generate
```

Expected: `pkg/openapi/server.gen.go` is regenerated (the `DuplicateProgramTemplate` interface method signature should NOT change since oapi-codegen passes request body via `r *http.Request`).

- [ ] **Step 3: Commit**

```bash
git add api/openapi/openapi.yaml api/pkg/openapi/server.gen.go
git commit -m "feat: add optional name field to duplicate template endpoint"
```

---

### Task 2: Update Go handler to use optional name from request body

**Files:**
- Modify: `api/internal/handler/program_template.go:189-266`

- [ ] **Step 1: Update DuplicateProgramTemplate handler to parse optional body**

In `api/internal/handler/program_template.go`, replace the `copyName` logic in `DuplicateProgramTemplate`:

After getting the `original` (line 207), replace the single line `copyName := ...` with:

```go
	// Parse optional request body for custom name
	copyName := strings.TrimSpace(original.Name) + " (copy)"
	if r.ContentLength > 0 {
		var req struct {
			Name string `json:"name,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			middleware.WriteValidationError(w, "Invalid request body", nil)
			return
		}
		if req.Name != "" {
			copyName = req.Name
		}
	}
```

- [ ] **Step 2: Run tests and lint**

Run from `api/`:
```bash
task check
```

Expected: All checks pass.

- [ ] **Step 3: Commit**

```bash
git add api/internal/handler/program_template.go
git commit -m "feat: support optional name in duplicate template handler"
```

---

### Task 3: Update frontend API client and hook

**Files:**
- Modify: `web/lib/api/programTemplates.ts:46-48`
- Modify: `web/lib/hooks/useProgramTemplates.ts:66-75`

- [ ] **Step 1: Update API client to accept optional name**

In `web/lib/api/programTemplates.ts`, change the `duplicate` method:

```typescript
  async duplicate(id: string, name?: string): Promise<ProgramTemplate> {
    const options = name ? { json: { name } } : undefined;
    return api.post(`program-templates/${id}/duplicate`, options).json<ProgramTemplate>();
  },
```

- [ ] **Step 2: Update the hook to accept `{ id, name? }`**

In `web/lib/hooks/useProgramTemplates.ts`, change `useDuplicateProgramTemplate`:

```typescript
export function useDuplicateProgramTemplate() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, name }: { id: string; name?: string }) =>
      programTemplatesApi.duplicate(id, name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['program-templates'] });
    },
  });
}
```

- [ ] **Step 3: Run lint**

Run from `web/`:
```bash
pnpm check
```

Expected: No errors.

- [ ] **Step 4: Commit**

```bash
git add web/lib/api/programTemplates.ts web/lib/hooks/useProgramTemplates.ts
git commit -m "feat: support optional name parameter in duplicate template API"
```

---

### Task 4: Add Duplicate and Archive confirmation dialogs to detail page

**Files:**
- Modify: `web/app/program-templates/[id]/page.tsx`

- [ ] **Step 1: Add state and dialog for Duplicate**

In `page.tsx`, add new state variables after the existing `deleteDialogOpen` state:

```typescript
const [duplicateDialogOpen, setDuplicateDialogOpen] = useState(false);
const [duplicateName, setDuplicateName] = useState('');
const [archiveDialogOpen, setArchiveDialogOpen] = useState(false);
```

Add imports at the top:

```typescript
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
```

- [ ] **Step 2: Update handlers**

Replace `handleDuplicate`:

```typescript
const openDuplicateDialog = () => {
  setDuplicateName(`${template.name} (copy)`);
  setDuplicateDialogOpen(true);
};

const handleDuplicate = async () => {
  const name = duplicateName.trim();
  await duplicateTemplate.mutateAsync({ id: templateId, name: name || undefined });
  setDuplicateDialogOpen(false);
  router.push('/program-templates');
};
```

Replace `handleArchiveToggle`:

```typescript
const handleArchiveToggle = async () => {
  if (isArchived) {
    await unarchiveTemplate.mutateAsync(templateId);
  } else {
    setArchiveDialogOpen(true);
  }
};

const handleArchive = async () => {
  await archiveTemplate.mutateAsync(templateId);
  setArchiveDialogOpen(false);
  router.push('/program-templates');
};
```

- [ ] **Step 3: Update Duplicate button onClick**

Change the Duplicate button's `onClick` from `handleDuplicate` to `openDuplicateDialog`.

- [ ] **Step 4: Add Duplicate dialog JSX**

After the existing `DeleteConfirmDialog`, add:

```tsx
<Dialog open={duplicateDialogOpen} onOpenChange={setDuplicateDialogOpen}>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Duplicate template</DialogTitle>
      <DialogDescription>
        Enter a name for the duplicated template.
      </DialogDescription>
    </DialogHeader>
    <div className="space-y-2">
      <Label htmlFor="duplicate-name">Name</Label>
      <Input
        id="duplicate-name"
        value={duplicateName}
        onChange={(e) => setDuplicateName(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === 'Enter' && duplicateName.trim()) {
            handleDuplicate();
          }
        }}
      />
    </div>
    <DialogFooter>
      <Button
        variant="outline"
        onClick={() => setDuplicateDialogOpen(false)}
      >
        Cancel
      </Button>
      <Button
        onClick={handleDuplicate}
        disabled={!duplicateName.trim() || duplicateTemplate.isPending}
      >
        Duplicate
      </Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

- [ ] **Step 5: Add Archive confirmation dialog JSX**

After the Duplicate dialog, add:

```tsx
<AlertDialog open={archiveDialogOpen} onOpenChange={setArchiveDialogOpen}>
  <AlertDialogContent>
    <AlertDialogHeader>
      <AlertDialogTitle>Archive template?</AlertDialogTitle>
      <AlertDialogDescription>
        This template will be hidden from the main list. You can unarchive it later.
      </AlertDialogDescription>
    </AlertDialogHeader>
    <AlertDialogFooter>
      <AlertDialogCancel>Cancel</AlertDialogCancel>
      <AlertDialogAction onClick={handleArchive}>
        Archive
      </AlertDialogAction>
    </AlertDialogFooter>
  </AlertDialogContent>
</AlertDialog>
```

- [ ] **Step 6: Run lint**

Run from `web/`:
```bash
pnpm check
```

Expected: No errors.

- [ ] **Step 7: Commit**

```bash
git add web/app/program-templates/[id]/page.tsx
git commit -m "feat: add confirmation dialogs for duplicate and archive template actions"
```
