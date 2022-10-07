<script lang="ts">
	import Footer from '../../../components/Footer.svelte';
	import Header from '../../../components/Header.svelte';
  import QuoteCard from '../../../components/QuoteCard.svelte';
</script>

<Header />

<div class=" h-full text-center p-4">
  <h1 class="text-3xl font-bold lg:text-5xl lg:font-black pb-10">Submitted Definitions</h1>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 justify-around">
    <div class="bg-white card grid justify-center shadow-2xl">
      <div class="grid p-3 text-xl content-between">
        <div id="forCategAndSymb" class="bg-blue-900 flex justify-between h-12">
          <div id="forCaterg" class=" font-normal self-center">
            &#128083; Artificial Intelligence
          </div>
          <div id="forSymbol" class="flex flex-wrap gap-3 items-center">
            <span title="accept" style='font-size:20px;'>&#128077;</span>
            <span title="don't accept" style='font-size:20px;'>&#128078;</span>
          </div>
        </div>
        <div id="forDefinition" class="bg-blue-100 font-light h-auto align-middle my-4">
          "Ab totam sint et natus velit qui nihil totam et laboriosam asperiores vel voluptatibus odio id ratione nesciunt. 
          Ut dignissimos voluptas ut quos tempore sit aliquam magnam est delectus aperiam quo quia suscipit in consequuntur sunt. 
          Hic nemo eligendi eum quam quae aut molestiae blanditiis et commodi unde?"
        </div>
        <div id="forAuthor" class="flex justify-end font-semibold bg-blue-300">
          Hamid Omar
        </div>
        <div id="forSource" class="bg-blue-700 font-normal text-base flex justify-end">
          From: 12 Rule for Life
        </div>
      </div>
    </div>
    <div class="card bg-white grid justify-center shadow-2xl">
      <div class="grid p-3 text-xl ceontent-between">
        <div id="forCategAndSymb" class=" flex justify-between h-12">
          <div id="forCaterg" class=" badge badge-secondary text-base mb-2 p-3">
            &#128083; Artificial Intelligence
          </div>
          <div id="forSymbol" class="flex flex-wrap gap-3 items-center">
            <span title="accept" class="text-base">&#128077;</span>
            <span title="don't accept" class="text-base">&#128078;</span>
          </div>
        </div>
        <div id="forDefinition" class=" font-light h-auto align-middle my-4">
          "Ab totam sint et natus velit qui nihil totam et laboriosam asperiores vel voluptatibus odio id ratione nesciunt. 
          Ut dignissimos voluptas ut quos tempore sit aliquam magnam est delectus aperiam quo quia suscipit in consequuntur sunt. 
          Hic nemo eligendi eum quam quae aut molestiae blanditiis et commodi unde?"
        </div>
        <div id="forAuthor" class="flex justify-end font-semibold">
          Hamid Omar
        </div>
        <div id="forSource" class=" font-normal text-base flex justify-end">
          From: 12 Rule for Life
        </div>
      </div>
    </div>
    <div class="card bg-white grid justify-center shadow-2xl">
      <div class="grid p-3 text-xl ceontent-between">
        <div id="forCategAndSymb" class=" flex justify-between h-12">
          <div id="forCaterg" class=" font-normal self-center bg-blue-300">
            &#128083; Artificial Intelligence
          </div>
          <div id="forSymbol" class="flex flex-wrap gap-3 items-center">
            <span title="accept" class="text-base">&#128077;</span>
            <span title="don't accept" class="text-base">&#128078;</span>
          </div>
        </div>
        <div id="forDefinition" class=" font-light h-auto align-middle my-4">
          "Ab totam sint et natus velit qui nihil totam et laboriosam asperiores vel voluptatibus odio id ratione nesciunt. 
          Ut dignissimos voluptas ut quos tempore sit aliquam magnam est delectus aperiam quo quia suscipit in consequuntur sunt. 
          Hic nemo eligendi eum quam quae aut molestiae blanditiis et commodi unde?"
        </div>
        <div id="forAuthor" class="flex justify-end font-semibold">
          Hamid Omar
        </div>
        <div id="forSource" class=" font-normal text-base flex justify-end">
          From: 12 Rule for Life
        </div>
      </div>
    </div>
    <div class="card bg-white grid justify-center shadow-2xl">
      <div class="grid p-3 text-xl ceontent-between">
        <div id="forCategAndSymb" class=" flex justify-between h-12">
          <div id="forCaterg" class=" font-normal self-center">
            &#128083; Artificial Intelligence
          </div>
          <div id="forSymbol" class="flex flex-wrap gap-3 items-center">
            <span title="accept" class="text-base">&#128077;</span>
            <span title="don't accept" class="text-base">&#128078;</span>
          </div>
        </div>
        <div id="forDefinition" class=" font-light h-auto align-middle my-4">
          "Ab totam sint et natus velit qui nihil totam et laboriosam asperiores vel voluptatibus odio id ratione nesciunt. 
          Ut dignissimos voluptas ut quos tempore sit aliquam magnam est delectus aperiam quo quia suscipit in consequuntur sunt. 
          Hic nemo eligendi eum quam quae aut molestiae blanditiis et commodi unde?"
        </div>
        <div id="forAuthor" class="flex justify-end font-semibold">
          Hamid Omar
        </div>
        <div id="forSource" class=" font-normal text-base flex justify-end">
          From: 12 Rule for Life
        </div>
      </div>
    </div>
    <QuoteCard
    category="Machine Intelligence"
    quote="Artificial Intelligence is […] the study of the computations that make it possible to perceive, reason, and act.
    Ab totam sint et natus velit qui nihil totam et laboriosam asperiores vel voluptatibus odio id ratione nesciunt. 
    Ut dignissimos voluptas ut quos tempore sit aliquam magnam est delectus aperiam quo quia suscipit in consequuntur sunt. 
    Hic nemo eligendi eum quam quae aut molestiae blanditiis et commodi unde?"
    author="Winston, P.H."
    publishing_date="1992"
    action="true"
  />

  <QuoteCard
    category="Machine Intelligence"
    quote="Artificial Intelligence is […] the study of the computations that make it possible to perceive, reason, and act."
    author="Winston, P.H."
    publishing_date="1992"
    action="true"
  />
  </div>
</div>

<Footer />

<style lang="postcss">
</style>
