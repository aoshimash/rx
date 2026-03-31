'use client';

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { DeleteConfirmDialog } from '@/components/ui/delete-confirm-dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import {
  useCreateFieldGroup,
  useDeleteFieldGroup,
  useFieldGroups,
  useUpdateFieldGroup,
} from '@/lib/hooks/useFieldGroups';
import type { FieldDef, FieldGroup } from '@/types/api';
import { Pencil, Plus, Trash2, X } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

// ============================================================================
// FieldDefListEditor — edits an array of FieldDef
// ============================================================================

function FieldDefListEditor({
  fields,
  onChange,
  label,
  isLogFields = false,
}: {
  fields: FieldDef[];
  onChange: (fields: FieldDef[]) => void;
  label: string;
  isLogFields?: boolean;
}) {
  const [newName, setNewName] = useState('');
  const [newType, setNewType] = useState<'text' | 'number' | 'select' | 'video'>('number');

  const addField = () => {
    const trimmed = newName.trim();
    if (trimmed && !fields.some((f) => f.name === trimmed)) {
      onChange([...fields, { name: trimmed, type: newType }]);
      setNewName('');
      setNewType('number');
    }
  };

  return (
    <div className="space-y-2">
      <Label className="text-xs font-semibold uppercase text-muted-foreground">{label}</Label>
      {fields.map((field, idx) => (
        <div key={field.name} className="flex items-start gap-2 rounded border p-2">
          <div className="flex-1 space-y-1">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium">{field.name}</span>
              <span className="text-xs text-muted-foreground">({field.type})</span>
            </div>
            <Input
              value={field.description ?? ''}
              onChange={(e) => {
                const next = [...fields];
                next[idx] = { ...field, description: e.target.value || undefined };
                onChange(next);
              }}
              placeholder="Description (e.g., unit, scale)"
              className="h-7 text-xs"
            />
          </div>
          <button
            type="button"
            onClick={() => onChange(fields.filter((_, i) => i !== idx))}
            className="mt-1 shrink-0 text-muted-foreground hover:text-foreground cursor-pointer"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        </div>
      ))}
      <div className="flex items-center gap-2">
        <Input
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              e.preventDefault();
              addField();
            }
          }}
          placeholder="Field name"
          className="h-8 w-40"
        />
        <Select
          value={newType}
          onValueChange={(v) => setNewType(v as 'text' | 'number' | 'select' | 'video')}
        >
          <SelectTrigger className="h-8 w-28">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="number">Number</SelectItem>
            <SelectItem value="text">Text</SelectItem>
            <SelectItem value="select">Select</SelectItem>
            {isLogFields && <SelectItem value="video">Video</SelectItem>}
          </SelectContent>
        </Select>
        <Button variant="ghost" size="sm" onClick={addField} disabled={!newName.trim()}>
          <Plus className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}

// ============================================================================
// FieldGroupForm — create or edit a single field group
// ============================================================================

interface FieldGroupFormProps {
  initial?: FieldGroup;
  onSave: (data: {
    name: string;
    description?: string;
    program_fields: FieldDef[];
    log_fields: FieldDef[];
  }) => void;
  onCancel: () => void;
  isSaving: boolean;
}

function FieldGroupForm({ initial, onSave, onCancel, isSaving }: FieldGroupFormProps) {
  const [name, setName] = useState(initial?.name ?? '');
  const [description, setDescription] = useState(initial?.description ?? '');
  const [programFields, setProgramFields] = useState<FieldDef[]>(initial?.program_fields ?? []);
  const [logFields, setLogFields] = useState<FieldDef[]>(initial?.log_fields ?? []);

  const handleSubmit = () => {
    if (!name.trim() || programFields.length === 0 || logFields.length === 0) return;
    onSave({
      name: name.trim(),
      description: description.trim() || undefined,
      program_fields: programFields,
      log_fields: logFields,
    });
  };

  return (
    <div className="space-y-4 rounded-lg border p-4">
      <div className="space-y-2">
        <Label htmlFor="fg-name">Name</Label>
        <Input
          id="fg-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g., Strength, Conditioning"
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="fg-desc">Description</Label>
        <Textarea
          id="fg-desc"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="What this field group is for (helps AI understand context)"
          rows={2}
        />
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <FieldDefListEditor
          fields={programFields}
          onChange={setProgramFields}
          label="Program Fields (Plan)"
        />
        <FieldDefListEditor
          fields={logFields}
          onChange={setLogFields}
          label="Log Fields (Record)"
          isLogFields
        />
      </div>

      <div className="flex gap-2 justify-end">
        <Button variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          disabled={
            isSaving || !name.trim() || programFields.length === 0 || logFields.length === 0
          }
        >
          {isSaving ? 'Saving...' : initial ? 'Update' : 'Create'}
        </Button>
      </div>
    </div>
  );
}

// ============================================================================
// FieldGroupCard — displays one field group with edit/delete
// ============================================================================

function FieldGroupCard({
  group,
  onEdit,
  onDelete,
}: {
  group: FieldGroup;
  onEdit: () => void;
  onDelete: () => void;
}) {
  const [deleteOpen, setDeleteOpen] = useState(false);

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-base">{group.name}</CardTitle>
          <div className="flex gap-1">
            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={onEdit}>
              <Pencil className="h-3.5 w-3.5" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-destructive"
              onClick={() => setDeleteOpen(true)}
            >
              <Trash2 className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>
        {group.description && <p className="text-sm text-muted-foreground">{group.description}</p>}
      </CardHeader>
      <CardContent>
        <Accordion type="single" collapsible>
          <AccordionItem value="fields" className="border-0">
            <AccordionTrigger className="py-1 text-xs text-muted-foreground hover:no-underline">
              {group.program_fields.length} program fields, {group.log_fields.length} log fields
            </AccordionTrigger>
            <AccordionContent>
              <div className="grid gap-3 md:grid-cols-2 pt-2">
                <div>
                  <p className="text-xs font-semibold uppercase text-muted-foreground mb-1">
                    Program Fields
                  </p>
                  <ul className="space-y-0.5">
                    {group.program_fields.map((f) => (
                      <li key={f.name} className="text-sm">
                        <span className="font-medium">{f.name}</span>
                        <span className="text-muted-foreground"> ({f.type})</span>
                        {f.description && (
                          <span className="text-muted-foreground"> — {f.description}</span>
                        )}
                      </li>
                    ))}
                  </ul>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase text-muted-foreground mb-1">
                    Log Fields
                  </p>
                  <ul className="space-y-0.5">
                    {group.log_fields.map((f) => (
                      <li key={f.name} className="text-sm">
                        <span className="font-medium">{f.name}</span>
                        <span className="text-muted-foreground"> ({f.type})</span>
                        {f.description && (
                          <span className="text-muted-foreground"> — {f.description}</span>
                        )}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </CardContent>

      <DeleteConfirmDialog
        open={deleteOpen}
        onOpenChange={setDeleteOpen}
        onConfirm={() => {
          onDelete();
          setDeleteOpen(false);
        }}
        title="Delete Field Group?"
        description={`This will permanently delete "${group.name}". Sessions using this group will lose their field definitions.`}
      />
    </Card>
  );
}

// ============================================================================
// FieldGroupEditor — main component for Settings page
// ============================================================================

export function FieldGroupEditor() {
  const { data, isLoading } = useFieldGroups();
  const createFieldGroup = useCreateFieldGroup();
  const updateFieldGroup = useUpdateFieldGroup();
  const deleteFieldGroup = useDeleteFieldGroup();

  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const groups = data?.data ?? [];

  const handleCreate = (formData: {
    name: string;
    description?: string;
    program_fields: FieldDef[];
    log_fields: FieldDef[];
  }) => {
    createFieldGroup.mutate(formData, {
      onSuccess: () => {
        setShowCreateForm(false);
        toast.success('Field group created');
      },
      onError: () => toast.error('Failed to create field group'),
    });
  };

  const handleUpdate = (
    id: string,
    formData: {
      name: string;
      description?: string;
      program_fields: FieldDef[];
      log_fields: FieldDef[];
    }
  ) => {
    updateFieldGroup.mutate(
      { id, data: formData },
      {
        onSuccess: () => {
          setEditingId(null);
          toast.success('Field group updated');
        },
        onError: () => toast.error('Failed to update field group'),
      }
    );
  };

  const handleDelete = (id: string) => {
    deleteFieldGroup.mutate(id, {
      onSuccess: () => toast.success('Field group deleted'),
      onError: () => toast.error('Failed to delete field group'),
    });
  };

  if (isLoading) {
    return <p className="text-muted-foreground">Loading field groups...</p>;
  }

  return (
    <div className="space-y-4">
      {groups.map((group) =>
        editingId === group.id ? (
          <FieldGroupForm
            key={group.id}
            initial={group}
            onSave={(data) => handleUpdate(group.id, data)}
            onCancel={() => setEditingId(null)}
            isSaving={updateFieldGroup.isPending}
          />
        ) : (
          <FieldGroupCard
            key={group.id}
            group={group}
            onEdit={() => setEditingId(group.id)}
            onDelete={() => handleDelete(group.id)}
          />
        )
      )}

      {groups.length === 0 && !showCreateForm && (
        <p className="text-sm text-muted-foreground">
          No field groups yet. Create one to define what fields your sessions use.
        </p>
      )}

      {showCreateForm ? (
        <FieldGroupForm
          onSave={handleCreate}
          onCancel={() => setShowCreateForm(false)}
          isSaving={createFieldGroup.isPending}
        />
      ) : (
        <Button variant="outline" onClick={() => setShowCreateForm(true)} className="w-full">
          <Plus className="h-4 w-4 mr-2" />
          New Field Group
        </Button>
      )}
    </div>
  );
}
