INSERT INTO prompt_versions (version, label, active) VALUES ('v1-default', 'Default tutor prompt', true)
ON CONFLICT (version) DO UPDATE SET label = EXCLUDED.label, active = true;

UPDATE prompt_versions SET active = (version = 'v1-default');

INSERT INTO experiments (key, name, variants, active) VALUES
  ('tutor_correction_style', 'Strict vs gentle corrections', ARRAY['strict', 'gentle'], true),
  ('tutor_tone', 'Professional vs friendly tutor', ARRAY['professional', 'friendly'], true),
  ('onboarding_length', 'Onboarding depth', ARRAY['full', 'lightweight'], true)
ON CONFLICT (key) DO NOTHING;

INSERT INTO feature_flags (key, enabled, description) VALUES
  ('livekit_voice_enabled', true, 'Enable LiveKit voice sessions'),
  ('post_session_llm_summary', true, 'Generate LLM post-session feedback'),
  ('analytics_debug_logging', false, 'Verbose analytics payload logging')
ON CONFLICT (key) DO NOTHING;
