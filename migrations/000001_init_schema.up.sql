-- 000001_init_schema.up.sql
BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    code_name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Permission Groups
CREATE TABLE permission_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Permission Group Items (many2many)
CREATE TABLE permission_group_items (
    permission_group_id UUID NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (permission_group_id, permission_id)
);

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_email ON users(email);

-- User Permissions
CREATE TABLE user_permissions (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission_group_id UUID NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, permission_group_id)
);

-- Verification Tokens
CREATE TABLE verification_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    code_hash CHAR(64) NOT NULL,
    expire_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_verification_tokens_user_id ON verification_tokens(user_id);
CREATE UNIQUE INDEX idx_verification_tokens_code_hash ON verification_tokens(code_hash);
CREATE INDEX idx_verification_tokens_type ON verification_tokens(type);
CREATE INDEX idx_verification_tokens_expire_at ON verification_tokens(expire_at);

-- Procedures
CREATE TABLE procedures (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES procedures(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    version INTEGER NOT NULL DEFAULT 1
);

-- Procedure Steps
CREATE TABLE procedure_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    procedure_id UUID NOT NULL REFERENCES procedures(id) ON DELETE CASCADE,
    index INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    is_optional BOOLEAN NOT NULL DEFAULT FALSE,
    step_type VARCHAR(20) NOT NULL,
    wait_time INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_procedure_steps_procedure_id ON procedure_steps(procedure_id);
CREATE INDEX idx_procedure_steps_step_type ON procedure_steps(step_type);

-- Experiments
CREATE TABLE experiments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    objective TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    type VARCHAR(20) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_by_id UUID NOT NULL REFERENCES users(id),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    procedure_id UUID REFERENCES procedures(id)
);

CREATE INDEX idx_experiments_status ON experiments(status);
CREATE INDEX idx_experiments_type ON experiments(type);
CREATE INDEX idx_experiments_created_by_id ON experiments(created_by_id);
CREATE INDEX idx_experiments_procedure_id ON experiments(procedure_id);

-- Experiment Results
CREATE TABLE experiment_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    experiment_id UUID NOT NULL UNIQUE REFERENCES experiments(id) ON DELETE CASCADE,
    outcome VARCHAR(20) NOT NULL,
    summary TEXT NOT NULL,
    outcome_reason TEXT NOT NULL,
    confidence_level VARCHAR(20) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Observations
CREATE TABLE observations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    observed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    title TEXT NOT NULL,
    notes TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    previous_observation_id UUID,
    experiment_id UUID NOT NULL,
    procedure_step_id UUID,

    CONSTRAINT observations_previous_observation_id_fkey FOREIGN KEY (previous_observation_id) REFERENCES observations(id) ON DELETE SET NULL,
    CONSTRAINT observations_experiment_id_fkey FOREIGN KEY (experiment_id) REFERENCES experiments(id) ON DELETE CASCADE,
    CONSTRAINT observations_procedure_step_id_fkey FOREIGN KEY (procedure_step_id) REFERENCES procedure_steps(id) ON DELETE SET NULL
);

CREATE INDEX idx_observations_observed_at ON observations(observed_at);
CREATE INDEX idx_observations_experiment_id ON observations(experiment_id);
CREATE INDEX idx_observations_procedure_step_id ON observations(procedure_step_id);

COMMIT;
