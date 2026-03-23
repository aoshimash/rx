ALTER TABLE program_templates
  ADD COLUMN source_template_id UUID REFERENCES program_templates(id) ON DELETE SET NULL;
