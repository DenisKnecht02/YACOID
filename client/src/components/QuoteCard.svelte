<script lang="ts">
    export let category: string;
    export let quote: string;
    export let author: string;
    export let publishing_date: string;
    export let source: string;
    export let action: string;
</script>

<div class="card w-[1/4] bg-base-100 shadow-xl">
    <div class="card-body">
        <div class="flex justify-between h-12 items-center">
          {#if action === "true"}
          <div class="badge badge-secondary text-base mb-2 p-3">{category}</div>
          <div id="forSymbol" class="flex flex-wrap gap-3">
            <span title="accept" class="text-base">&#128077;</span>
            <span title="don't accept" class="text-base">&#128078;</span>
          </div>
            
          {:else}
          <div class="badge badge-secondary text-base mb-2 p-3">{category}</div>
          {/if}
        </div>
        <p class="flex italic text-xl text-justify items-center">"{quote}"</p>
        <div>
          <p class="flex justify-end items-end font-bold">- {author} ({publishing_date})</p>
          <p class="flex justify-end items-end italic">from: {source}</p>
        </div>
        
    </div>
</div>

<style lang="postcss">
    
</style>