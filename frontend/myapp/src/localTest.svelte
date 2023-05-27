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
    let rawData = [];
    let board = "";
    let reason = "";
    let result = "";
    let testName = "";
    let startDate = "";
    let endDate = "";
    let orderBy = "";
    let selectResult = [
        { value: "", name: "All" },
        { value: "Pass", name: "Pass" },
        { value: "Fail", name: "Fail" },
    ];

    const fetchData = async () => {
        let response = await fetch(
            `http://10.240.102.16:8082/localTest?startDate=${startDate}&endDate=${endDate}&board=${board}&reason=${reason}&testName=${testName}&result=${result}&orderBy=${orderBy}`
        );
        rawData = await response.json();
    };

    const fetch_data_2 = async () => {
        rawData = [];
        let response = await fetch(
            `http://10.240.102.16:8082/localTest?startDate=${startDate}&endDate=${endDate}&board=${board}&reason=${reason}&testName=${testName}&result=${result}&orderBy=${orderBy}`
        );
        rawData = await response.json();
    };

    const reset = () => {
        rawData = [];
    };

    onMount(async () => {
        await fetchData();
    });

    $: filteredItems = rawData.filter(
        (data) =>
            data.name.toLowerCase().indexOf(searchTerm.toLowerCase()) != -1 &&
            data.board.toLowerCase().indexOf(searchBoard.toLowerCase()) != -1 &&
            data.passOrFail.toLowerCase().indexOf(searchResult.toLowerCase()) !=
                -1
    );
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
                <div style="padding-top: 1.5rem;">
                    <button
                        style="background-color:wheat;"
                        on:click={() => {
                            reset();
                            fetch_data_2();
                        }}>Search</button
                    >
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
            placeholder="Search by result"
            hoverable={true}
            bind:inputValue={searchResult}
            striped={true}
            divClass="tableSearch"
        />
    </span>

    <Table hoverable={true} striped={true}>
        <TableHead theadClass="colorTable">
            <TableHeadCell class="dataInTable">ID</TableHeadCell>
            <TableHeadCell class="dataInTable">Time</TableHeadCell>
            <TableHeadCell class="dataInTable">Tester</TableHeadCell>
            <TableHeadCell class="dataInTable">Name</TableHeadCell>
            <TableHeadCell class="dataInTable">Board</TableHeadCell>
            <TableHeadCell class="dataInTable">Model</TableHeadCell>
            <TableHeadCell class="dataInTable">Version</TableHeadCell>
            <TableHeadCell class="dataInTable">Log_Path</TableHeadCell>
            <TableHeadCell class="dataInTable">Result</TableHeadCell>
            <TableHeadCell class="dataInTable">Reason</TableHeadCell>
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
                            >{post.tester}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.name}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.board}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.model}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.version}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable">
                            <a
                                style="color: white;"
                                href="http://10.240.102.16:8082/logDB/{post.logPath}"
                                >{post.logPath}</a
                            >
                        </TableBodyCell>
                        <TableBodyCell tdClass="dataInTable"
                            >{post.passOrFail}</TableBodyCell
                        >
                        <TableBodyCell tdClass="dataInTable"
                            >{post.reason}</TableBodyCell
                        >
                    </TableBodyRow>
                {/each}
            {/if}
        </TableBody>
    </Table>
</main>

<style>
    .line {
        width: 100vw;
        height: 10vh;
        background: #242424;
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
        left: 18px;
        top: 27px;
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
        height: 50vh;
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

    /* .menu-list ul {
        list-style: none;
        margin-left: 70px;
        padding: 0;
    }
    .menu-list ul > li {
        display: block;
        width: 100px;
    } */

    input#menu {
        display: none;
    }

    span {
        color: rgb(250, 33, 33);
        font-weight: bold;
    }
</style>
