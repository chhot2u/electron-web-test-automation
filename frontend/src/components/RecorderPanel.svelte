<script lang="ts">
  import { StartRecording, StopRecording, CreateRecordedFlow } from '../../wailsjs/go/main/App';
  import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
  import { isRecording, recordingSteps } from '../lib/store';
  import type { RecordedStep } from '../lib/types';
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';

  const dispatch = createEventDispatcher();

  let flowName = '';
  let flowDescription = '';
  let originUrl = '';
  let errorMessage = '';
  let saving = false;
  let starting = false;

  async function toggleRecording() {
    errorMessage = '';
    if ($isRecording) {
      await stopRecording();
    } else {
      await startRecording();
    }
  }

  async function startRecording() {
    if (!originUrl.trim()) {
      errorMessage = 'URL is required to start recording';
      return;
    }
    starting = true;
    try {
      recordingSteps.set([]);
      await StartRecording(originUrl);
      isRecording.set(true);
    } catch (err: any) {
      errorMessage = err?.message || String(err);
    } finally {
      starting = false;
    }
  }

  async function stopRecording() {
    try {
      const steps = await StopRecording();
      isRecording.set(false);
      recordingSteps.set(steps || []);
    } catch (err: any) {
      errorMessage = err?.message || String(err);
    }
  }

  function handleRecorderStep(step: RecordedStep) {
    recordingSteps.update(list => [...list, step]);
  }

  async function saveFlow() {
    if (!flowName || $recordingSteps.length === 0) return;
    saving = true;
    try {
      errorMessage = '';
      await CreateRecordedFlow(flowName, flowDescription, originUrl, $recordingSteps);
      dispatch('saved');
      flowName = '';
      flowDescription = '';
      recordingSteps.set([]);
    } catch (err: any) {
      errorMessage = err?.message || String(err);
    } finally {
      saving = false;
    }
  }

  onMount(() => {
    EventsOn('recorder:step', handleRecorderStep);
  });

  onDestroy(() => {
    EventsOff('recorder:step');
  });
</script>

<div class="panel">
  <div class="panel-header">
    <h2>Live Recorder</h2>
    <button
      class="btn-primary"
      class:btn-danger={$isRecording}
      disabled={starting || (!$isRecording && !originUrl.trim())}
      on:click={toggleRecording}
    >
      {#if starting}
        Launching...
      {:else if $isRecording}
        ⏹ Stop Recording
      {:else}
        ⏺ Start Recording
      {/if}
    </button>
  </div>

  <div class="panel-body">
    <div class="form-group">
      <label for="origin-url">URL</label>
      <input
        id="origin-url"
        bind:value={originUrl}
        placeholder="https://example.com"
        disabled={$isRecording}
      />
    </div>
    <div class="form-group">
      <label for="flow-name">Flow Name</label>
      <input id="flow-name" bind:value={flowName} placeholder="Checkout flow" />
    </div>
    <div class="form-group">
      <label for="flow-desc">Description</label>
      <input id="flow-desc" bind:value={flowDescription} placeholder="Optional" />
    </div>

    {#if $isRecording}
      <div class="recording-indicator">
        <span class="pulse"></span> Recording — interact with the browser window
      </div>
    {/if}

    <div class="steps">
      <h4>Recorded Steps ({$recordingSteps.length})</h4>
      {#if $recordingSteps.length === 0}
        <div class="empty">No steps recorded yet.</div>
      {:else}
        <ul>
          {#each $recordingSteps as step}
            <li>
              <strong>{step.action}</strong>
              {#if step.selector} <span class="muted">{step.selector}</span> {/if}
              {#if step.value} <span class="value">= {step.value}</span> {/if}
            </li>
          {/each}
        </ul>
      {/if}
    </div>
  </div>

  {#if errorMessage}
    <div class="error-banner">{errorMessage}</div>
  {/if}

  <div class="panel-footer">
    <button class="btn-secondary" on:click={() => dispatch('close')}>Close</button>
    <button
      class="btn-primary"
      disabled={!flowName || $recordingSteps.length === 0 || saving || $isRecording}
      on:click={saveFlow}
    >
      {saving ? 'Saving...' : 'Save Recorded Flow'}
    </button>
  </div>
</div>

<style>
  .panel {
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 16px;
  }
  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .panel-body {
    margin-top: 12px;
  }
  .recording-indicator {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    margin-top: 12px;
    background: rgba(239, 68, 68, 0.08);
    border: 1px solid rgba(239, 68, 68, 0.2);
    border-radius: 8px;
    color: var(--danger, #ef4444);
    font-size: 13px;
    font-weight: 500;
  }
  .pulse {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--danger, #ef4444);
    animation: pulse-anim 1.2s ease-in-out infinite;
  }
  @keyframes pulse-anim {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.3; }
  }
  .steps {
    margin-top: 16px;
  }
  .steps ul {
    list-style: none;
    padding: 0;
  }
  .steps li {
    padding: 6px 0;
    border-bottom: 1px solid var(--border);
  }
  .empty {
    font-size: 12px;
    color: var(--text-muted);
  }
  .muted {
    color: var(--text-muted);
    font-size: 11px;
  }
  .value {
    color: var(--text-muted);
    font-size: 11px;
    font-style: italic;
  }
  .panel-footer {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
  .btn-danger {
    background: var(--danger, #ef4444) !important;
    border-color: var(--danger, #ef4444) !important;
  }
</style>
