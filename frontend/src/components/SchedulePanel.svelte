<script lang="ts">
  import { onMount } from 'svelte';
  import { schedules, recordedFlows } from '../lib/store';
  import type { ProxyConfig } from '../lib/types';
  import { ListSchedules, CreateSchedule, ToggleSchedule, DeleteSchedule, ListRecordedFlows } from '../../wailsjs/go/main/App';

  let name = '';
  let cronExpr = '';
  let selectedFlowId = '';
  let url = '';
  let priority = 5;
  let headless = true;
  let tagsInput = '';
  let errorMessage = '';
  let creating = false;

  onMount(async () => {
    await refresh();
    await refreshFlows();
  });

  async function refresh() {
    try {
      errorMessage = '';
      const list = await ListSchedules();
      schedules.set(list || []);
    } catch (err: any) {
      errorMessage = `Failed to load schedules: ${err?.message || err}`;
    }
  }

  async function refreshFlows() {
    try {
      const flows = await ListRecordedFlows();
      recordedFlows.set(flows || []);
    } catch (_) {}
  }

  async function create() {
    if (!name || !cronExpr || !selectedFlowId || !url) return;
    creating = true;
    const proxy: ProxyConfig = { server: '' };
    const tags = tagsInput
      .split(',')
      .map(t => t.trim())
      .filter(t => t.length > 0);
    try {
      errorMessage = '';
      await CreateSchedule(name, cronExpr, selectedFlowId, url, proxy, priority, headless, tags);
      name = '';
      cronExpr = '';
      selectedFlowId = '';
      url = '';
      tagsInput = '';
      await refresh();
    } catch (err: any) {
      errorMessage = err?.message || String(err);
    } finally {
      creating = false;
    }
  }

  async function toggle(id: string, enabled: boolean) {
    try {
      errorMessage = '';
      await ToggleSchedule(id, !enabled);
      await refresh();
    } catch (err: any) {
      errorMessage = `Failed to toggle: ${err?.message || err}`;
    }
  }

  async function remove(id: string) {
    try {
      errorMessage = '';
      await DeleteSchedule(id);
      await refresh();
    } catch (err: any) {
      errorMessage = `Failed to delete: ${err?.message || err}`;
    }
  }

  function flowName(flowId: string): string {
    const flow = $recordedFlows.find(f => f.id === flowId);
    return flow?.name ?? flowId.slice(0, 8);
  }

  function formatTime(t: string | undefined): string {
    if (!t) return '—';
    try {
      return new Date(t).toLocaleString();
    } catch {
      return t;
    }
  }
</script>

<div class="schedule-panel">
  <h3>Schedules</h3>
  {#if errorMessage}
    <div class="error-text">{errorMessage}</div>
  {/if}

  <div class="create-form">
    <div class="form-row">
      <div class="form-group">
        <label for="sched-name">Name</label>
        <input id="sched-name" bind:value={name} placeholder="Daily scrape" />
      </div>
      <div class="form-group">
        <label for="sched-cron">Cron Expression</label>
        <input id="sched-cron" bind:value={cronExpr} placeholder="0 */6 * * *" />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group">
        <label for="sched-flow">Flow</label>
        <select id="sched-flow" bind:value={selectedFlowId}>
          <option value="">Select flow...</option>
          {#each $recordedFlows as flow}
            <option value={flow.id}>{flow.name}</option>
          {/each}
        </select>
      </div>
      <div class="form-group">
        <label for="sched-url">URL</label>
        <input id="sched-url" bind:value={url} placeholder="https://example.com" />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group">
        <label for="sched-priority">Priority</label>
        <select id="sched-priority" bind:value={priority}>
          <option value={1}>Low</option>
          <option value={5}>Normal</option>
          <option value={10}>High</option>
        </select>
      </div>
      <div class="form-group">
        <label for="sched-headless">Headless</label>
        <label class="checkbox">
          <input id="sched-headless" type="checkbox" bind:checked={headless} />
          Run headless
        </label>
      </div>
      <div class="form-group">
        <label for="sched-tags">Tags</label>
        <input id="sched-tags" bind:value={tagsInput} placeholder="tag1, tag2" />
      </div>
    </div>
    <div class="form-actions">
      <button class="btn-primary btn-sm" on:click={create} disabled={!name || !cronExpr || !selectedFlowId || !url || creating}>
        {creating ? 'Creating...' : 'Create Schedule'}
      </button>
    </div>
  </div>

  <div class="schedule-list">
    {#each $schedules as sched (sched.id)}
      <div class="schedule-item">
        <div class="schedule-info">
          <div class="schedule-top">
            <span class="schedule-name">{sched.name}</span>
            <span class="badge" class:badge-running={sched.enabled} class:badge-cancelled={!sched.enabled}>
              {sched.enabled ? 'enabled' : 'disabled'}
            </span>
          </div>
          <div class="schedule-meta">
            <span class="font-mono">{sched.cronExpr}</span>
            <span>| Flow: {flowName(sched.flowId)}</span>
          </div>
          <div class="schedule-times">
            <span>Last: {formatTime(sched.lastRunAt)}</span>
            <span>Next: {formatTime(sched.nextRunAt)}</span>
          </div>
        </div>
        <div class="schedule-actions">
          <button class="btn-secondary btn-sm" on:click={() => toggle(sched.id, sched.enabled)}>
            {sched.enabled ? 'Disable' : 'Enable'}
          </button>
          <button class="btn-danger btn-sm" on:click={() => remove(sched.id)}>Del</button>
        </div>
      </div>
    {:else}
      <p class="text-muted empty-msg">No schedules configured.</p>
    {/each}
  </div>
</div>

<style>
  .schedule-panel {
    padding: 16px;
  }
  .schedule-panel h3 {
    font-size: 16px;
    margin-bottom: 12px;
  }
  .create-form {
    background: var(--bg-secondary);
    padding: 12px;
    border-radius: var(--radius);
    border: 1px solid var(--border);
    margin-bottom: 16px;
  }
  .form-row {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
  }
  .form-group {
    flex: 1;
  }
  .form-group label {
    display: block;
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
    margin-bottom: 4px;
  }
  .form-group input,
  .form-group select {
    width: 100%;
  }
  .form-actions {
    display: flex;
    justify-content: flex-end;
  }
  .checkbox {
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
  }
  .checkbox input[type="checkbox"] {
    width: auto;
    padding: 0;
  }
  .schedule-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  .schedule-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 12px;
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: var(--radius);
  }
  .schedule-top {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .schedule-name {
    font-size: 13px;
    font-weight: 600;
  }
  .schedule-meta {
    font-size: 11px;
    color: var(--text-muted);
    margin-top: 2px;
  }
  .schedule-times {
    font-size: 11px;
    color: var(--text-muted);
    margin-top: 2px;
    display: flex;
    gap: 12px;
  }
  .schedule-actions {
    display: flex;
    gap: 6px;
    flex-shrink: 0;
  }
  .error-text {
    color: var(--danger, #ef4444);
    font-size: 12px;
    margin-bottom: 8px;
  }
  .empty-msg {
    text-align: center;
    padding: 20px;
  }
</style>
