DROP INDEX IF EXISTS idx_program_templates_active_name;

ALTER TABLE program_templates
  ADD COLUMN source_template_id UUID REFERENCES program_templates(id) ON DELETE SET NULL;
