<script lang="ts">
  import { onMount } from 'svelte';

  let analyzing = false;
  let result = '';
  let error = '';
  let visible = false;

  onMount(() => {
    // @ts-ignore - Wails runtime will be available
    window.runtime.EventsOn("analysis-result", (result: string) => {
      analyzing = false;
      visible = true;
      error = '';
      result = result;
    });
  });

  function closeWindow() {
    visible = false;
    result = '';
    error = '';
  }
</script>

{#if visible}
  <main class="result-window" class:analyzing>
    <div class="title-bar">
      <span>Analysis Result</span>
      <button class="close-btn" on:click={closeWindow}>×</button>
    </div>

    {#if analyzing}
      <div class="loading">
        Analyzing screenshot...
      </div>
    {:else if error}
      <div class="error">
        {error}
      </div>
    {:else if result}
      <div class="result">
        {result}
      </div>
    {/if}
  </main>
{/if}

<style>
  .result-window {
    position: fixed;
    top: 20px;
    right: 20px;
    width: 300px;
    max-height: 400px;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    overflow-y: auto;
  }

  .title-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 15px;
    background: #f5f5f5;
    border-top-left-radius: 8px;
    border-top-right-radius: 8px;
    border-bottom: 1px solid #eee;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 20px;
    cursor: pointer;
    padding: 0 5px;
  }

  .close-btn:hover {
    color: #f44336;
  }

  .analyzing {
    background: rgba(0, 0, 0, 0.8);
    color: white;
  }

  .loading {
    text-align: center;
    padding: 20px;
  }

  .error {
    color: #f44336;
    padding: 10px 15px;
  }

  .result {
    white-space: pre-wrap;
    line-height: 1.5;
    padding: 15px;
  }
</style>