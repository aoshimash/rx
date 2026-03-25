-- Remove template versioning (source_template_id) and enforce name uniqueness on active templates

ALTER TABLE program_templates DROP COLUMN source_template_id;

CREATE UNIQUE INDEX idx_program_templates_active_name
  ON program_templates (name) WHERE archived_at IS NULL;
