<script>
    import { onMount } from "svelte";
    import {
        Table,
        TableBody,
        TableBodyCell,
        TableBodyRow,
        TableHead,
        TableHeadCell,
        Input,
        Label,
        Select,
        TableSearch,
    } from "flowbite-svelte";

    let searchTerm = "";
    let searchBoard = "";
    let searchResult = "";
    let searchReason = "";

    let board = "";
    let reason = "";
    let result = "";
    let testName = "";
    let startDate = "";
    let endDate = "";
    let rawData = [];
    let orderBy = "";
    let search = false;

    let selectResult = [
        { value: "", name: "All" },
        { value: "Pass", name: "Pass" },
        { value: "Fail", name: "Fail" },
    ];

    let selectOrder = [
        { value: "Id", name: "ID" },
        { value: "Board", name: "Board" },
        { value: "Model", name: "Model" },
        { value: "Reason", name: "Reason" },
    ];

    const fetchData = async () => {
        rawData = [];
        let response = await fetch(
            `http://10.240.102.16:8082/stainlessSearch?startDate=${startDate}&endDate=${endDate}&board=${board}&reason=${reason}&testName=${testName}&result=${result}&orderBy=${orderBy}`
        );
        rawData = await response.json();
    };

    const reset = () => {
        rawData = [];
    };

    $: filteredItems =
        rawData?.filter((data) => {
            return (
                data.testName
                    .toLowerCase()
                    .includes(searchTerm.toLowerCase()) &&
                data.board.toLowerCase().includes(searchBoard.toLowerCase()) &&
                data.status
                    .toLowerCase()
                    .includes(searchResult.toLowerCase()) &&
                data.reason.toLowerCase().includes(searchReason.toLowerCase())
            );
        }) ?? [];
</script>

<main>
    <nav class="nav-box">
        <input type="checkbox" id="menu" />
        <label for="menu" class="line">
            <div class="menu" />
        </label>

        <div class="menu-list" style="text-align: left;">
            <div style="height: 10vh;" />
            <span class="grid gap-6 mb-6 md:grid-cols-4">
                <div>
                    <Label for="case"
                        ><span style="font-size: 1rem;">Case</span></Label
                    >
                    <Input
                        type="text"
                        id="case"
                        placeholder="BatteryCharging"
                        bind:value={testName}
                        size="lg"
                        required
                    />
                </div>
                <div>
                    <Label for="board" color="red"
                        ><span style="font-size: 1rem;">Board</span></Label
                    >
                    <Input
                        type="text"
                        id="board"
                        placeholder="volteer"
                        bind:value={board}
                        size="lg"
                        required
                    />
                </div>
                <div>
                    <Label for="reason" color="red"
                        ><span style="font-size: 1rem;">Reason</span></Label
                    >
                    <Input
                        type="text"
                        id="board"
                        placeholder="Failed to reboot to ..."
                        bind:value={reason}
                        size="lg"
                        required
                    />
                </div>
            </span>
            <span class="grid gap-6 mb-6 md:grid-cols-4">
                <div>
                    <Label for="case"
                        ><span style="font-size: 1rem;">Start date</span></Label
                    >
                    <Input
                        type="date"
                        id="case"
                        bind:value={startDate}
                        size="lg"
                        required
                    />
                </div>
                <div>
                    <Label for="case"
                        ><span style="font-size: 1rem;">End date</span></Label
                    >
                    <Input
                        type="date"
                        id="case"
                        bind:value={endDate}
                        size="lg"
                        required
                    />
                </div>
                <div>
                    <Label
                        ><span style="font-size: 1rem;">Result</span>
                        <Select
                            class="mt-2"
                            items={selectResult}
                            bind:value={result}
                        />
                    </Label>
                    <div />
                </div></span
            >
            <span class="grid gap-6 mb-6 md:grid-cols-4">
                <div>
                    <Label
                        ><span style="font-size: 1rem;">Order By</span>
                        <Select
                            class="mt-2"
                            items={selectOrder}
                            bind:value={orderBy}
                        />
                    </Label>
                </div>
                <div style="padding-top: 1.5rem;">
                    {#if search == false}
                        <button
                            on:click={() => {
                                reset();
                                fetchData();
                                search = true;
                            }}>Search</button
                        >
                    {/if}
                </div>
            </span>
        </div>
    </nav>

    <div style="height: 10vh;" />

    <span style="display: flex;">
        <TableSearch
            placeholder="Search by case name"
            hoverable={true}
            bind:inputValue={searchTerm}
            striped={true}
            divClass="tableSearch"
        />
        <TableSearch
            placeholder="Search by board"
            hoverable={true}
            bind:inputValue={searchBoard}
            striped={true}
            divClass="tableSearch"
        />
        <TableSearch
            placeholder="Search by status"
            hoverable={true}
            bind:inputValue={searchResult}
            striped={true}
            divClass="tableSearch"
        />
        <TableSearch
            placeholder="Search by reason"
            hoverable={true}
            bind:inputValue={searchReason}
            striped={true}
            divClass="tableSearch"
        />
    </span>

    <Table hoverable={true} striped={true}>
        <TableHead theadClass="colorTable">
            <TableHeadCell class="dataInTable">ID</TableHeadCell>
            <TableHeadCell class="dataInTable">Time</TableHeadCell>
            <TableHeadCell class="dataInTable">Duration</TableHeadCell>
            <TableHeadCell class="dataInTable">Suite</TableHeadCell>
            <TableHeadCell class="dataInTable">Board</TableHeadCell>
            <TableHeadCell class="dataInTable">Model</TableHeadCell>
            <TableHeadCell class="dataInTable">Build Version</TableHeadCell>
            <TableHeadCell class="dataInTable">Host</TableHeadCell>
            <TableHeadCell class="dataInTable">Test Name</TableHeadCell>
            <TableHeadCell class="dataInTable">Status</TableHeadCell>
            <TableHeadCell class="dataInTable">Reason</TableHeadCell>
            <TableHeadCell class="dataInTable">RO Version</TableHeadCell>
            <TableHeadCell class="dataInTable">RW Version</TableHeadCell>
        </TableHead>

        <TableBody>
            {#if filteredItems != null}
                {#each filteredItems as post (post.id)}
                    <TableBodyRow
                        color="custom"
                        class="odd:bg-gray-500 even:bg-gray-400 odd:dark:bg-gray-500 even:dark:bg-purple-400 hover:bg-purple-400"
                    >
                        <TableBodyCell tdClass="idInTable"
                            >{post.id}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.time.slice(5, 10)}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.duration}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.suite}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.board}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.model}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.buildVersion}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.host}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.testName}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.status}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.reason}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.firmwareROVersion}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.firmwareRWVersion}</TableBodyCell
                        >
                    </TableBodyRow>
                {/each}
                {(search = false)}
            {/if}
        </TableBody>
    </Table>
</main>

<style>
    .line {
        width: 100vw;
        height: 10vh;
        background: #87857b;
        cursor: pointer;
        display: block;
        padding: 16px;
        position: fixed;
        z-index: 2;
    }
    .line .menu,
    .line .menu::before,
    .line .menu::after {
        background: #fffefe;
        content: "";
        display: block;
        height: 3px;
        position: absolute;
        transition: background ease 0.3s, top ease 0.3s 0.3s,
            transform ease 0.3s;
        width: 3rem;
    }
    .line .menu {
        left: 3vw;
        top: 5vh;
    }

    .line .menu::before {
        top: -6px;
    }

    .line .menu::after {
        top: 6px;
    }

    #menu:checked + .line .menu {
        background: transparent;
    }

    #menu:checked + .line .menu::before {
        transform: rotate(45deg);
    }

    #menu:checked + .line .menu::after {
        transform: rotate(-45deg);
    }

    #menu:checked + .line .menu::before,
    #menu:checked + .line .menu::after {
        top: 0;
        transition: top ease 0.3s, transform ease 0.3s 0.3s;
    }

    #menu:checked ~ .menu-list {
        height: 100vh;
        background-color: #000000;
    }

    .menu-list {
        width: 100vw;
        height: 10vh;
        background: #242424;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        padding-top: 60px;
        position: fixed;
        z-index: 1;
        transition: all 0.6s;
        overflow: auto;
    }

    input#menu {
        display: none;
    }

    span {
        color: rgb(250, 33, 33);
        font-weight: bold;
    }
</style>
