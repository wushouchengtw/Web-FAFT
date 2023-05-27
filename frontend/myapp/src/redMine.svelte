<script>
  import {
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    TableSearch,
  } from "flowbite-svelte";
  import { onMount } from "svelte";
  let searchTerm = "";
  let rawData = [];

  onMount(async () => {
    let response = await fetch("http://10.240.102.16:8082/getTicketList");
    rawData = await response.json();
  });

  $: filteredItems = rawData.filter(
    (data) =>
      data.caseName.toLowerCase().indexOf(searchTerm.toLowerCase()) != -1
  );
</script>

<main>
  <TableSearch
    placeholder="Search by case name"
    hoverable={true}
    bind:inputValue={searchTerm}
    striped={true}
  >
    <TableHead theadClass="colorTable">
      <TableHeadCell class="dataInTable">Case</TableHeadCell>
      <TableHeadCell class="dataInTable">Owner</TableHeadCell>
      <TableHeadCell class="dataInTable">Tikcet</TableHeadCell>
    </TableHead>
    <TableBody>
      {#if filteredItems != null}
        {#each filteredItems as post (post.id)}
          <TableBodyRow
            color="custom"
            class="odd:bg-gray-500 even:bg-gray-400 odd:dark:bg-gray-500 even:dark:bg-purple-400 hover:bg-purple-400"
          >
            <TableBodyCell tdClass="dataInTable">{post.caseName}</TableBodyCell>
            <TableBodyCell tdClass="dataInTable">{post.owner}</TableBodyCell>
            <TableBodyCell tdClass="dataInTable">
              <a
                href="https://partnerissuetracker.corp.google.com/issues/{post.ticket}"
                style="color:blue;"
              >
                {post.ticket}
              </a>
            </TableBodyCell>
          </TableBodyRow>
        {/each}
      {/if}
    </TableBody>
  </TableSearch>
</main>
