<script>
    import { Input, Label } from "flowbite-svelte";
    import { scaleLinear } from "d3-scale";
    import {
        Table,
        TableBody,
        TableBodyCell,
        TableBodyRow,
        TableHead,
        TableHeadCell,
        Select,
        Progressbar,
    } from "flowbite-svelte";

    let testName = "";
    let data = [];
    let search = false;
    let owner = "";
    let dataOwner = [];

    let ownerList = [
        { value: "Arthur", name: "Arthur" },
        { value: "Carlos", name: "Carlos" },
        { value: "Kevin", name: "Kevin" },
        { value: "Peggy", name: "Peggy" },
        { value: "Yu-Chen", name: "Yu-Chen" },
        { value: "Wayne", name: "Wayne" },
    ];

    const fetchData = async () => {
        console.log(`${testName}`);
        let response = await fetch(
            `http://10.240.102.16:8082/passingRate?testName=${testName}`
        );
        data = await response.json();
        console.log(data);
        let ps = data.map((row) => {
            return {
                x: 28 - row.week * 7,
                y: row.passingRate * 100,
            };
        });
        points = ps;
    };

    const fetchOwner = async () => {
        console.log(`${owner}`);
        let response = await fetch(
            `http://10.240.102.16:8082/ownerPassingRate?owner=${owner}`
        );
        dataOwner = await response.json();
    };

    const reset = () => {
        data = [];
    };

    const resetOwner = () => {
        dataOwner = [];
    };

    // import points from "./data.js";
    const yTicks = [0, 50, 60, 70, 80, 90, 100];
    const xTicks = [28, 21, 14, 7];
    const padding = { top: 20, right: 15, bottom: 20, left: 25 };
    let width = 500;
    let height = 200;
    $: xScale = scaleLinear()
        .domain([minX, maxX])
        .range([padding.left, width - padding.right]);

    $: yScale = scaleLinear()
        .domain([Math.min.apply(null, yTicks), Math.max.apply(null, yTicks)])
        .range([height - padding.bottom, padding.top]);

    let points = [
        { x: 28, y: 0 },
        { x: 21, y: 0 },
        { x: 14, y: 0 },
        // { x: 7, y: 0 },
    ];
    points.push({ x: 7, y: 0 });

    $: minX = points[0].x;
    $: maxX = points[points.length - 1].x;
    $: path = `M${points
        .map((p) => `${xScale(p.x)},${yScale(p.y)}`)
        .join("L")}`;
    $: area = `${path}L${xScale(maxX)},${yScale(0)}L${xScale(minX)},${yScale(
        0
    )}Z`;

    function formatMobile(tick) {
        return (tick / 7).toString().slice(-2);
    }
</script>

<main>
    <div style="width: 100%;">
        <div style="padding-top: 3vh;">
            <div>
                <span class="grid gap-6 mb-6 md:grid-cols-4">
                    <div />
                    <div>
                        <Input
                            type="text"
                            id="case"
                            placeholder="BatteryCharging"
                            bind:value={testName}
                            size="lg"
                            required
                        />
                    </div>
                    <button
                        on:click={() => {
                            reset();
                            fetchData();
                            search = true;
                        }}>Search</button
                    >
                </span>
            </div>
            <br />
            <br />
            <div style="text-align: center; width: 100vw; font-size:5rem">
                Passing rate for the last 4 weeks
            </div>
            <br />
            <br />

            <Table hoverable={true} striped={true}>
                <TableHead theadClass="colorTable">
                    <TableHeadCell class="dataInTable"
                        >Time period</TableHeadCell
                    >
                    <TableHeadCell class="dataInTable">Total runs</TableHeadCell
                    >
                    <TableHeadCell class="dataInTable">Pass</TableHeadCell>
                    <TableHeadCell class="dataInTable">Fail times</TableHeadCell
                    >
                    <TableHeadCell class="dataInTable"
                        >Passing rate</TableHeadCell
                    >
                </TableHead>
                <TableBody>
                    {#if data != null}
                        {#each data as post (post.week)}
                            <TableBodyRow
                                color="custom"
                                class="odd:bg-gray-500 even:bg-gray-400 odd:dark:bg-gray-500 even:dark:bg-purple-400 hover:bg-purple-400"
                            >
                                <TableBodyCell tdClass="idInTable"
                                    >{post.startDate}-{post.endDate}</TableBodyCell
                                >
                                <TableBodyCell tdClass="idInTable"
                                    >{post.totalRun}</TableBodyCell
                                >
                                <TableBodyCell tdClass="idInTable"
                                    >{post.pass}</TableBodyCell
                                >
                                <TableBodyCell tdClass="idInTable"
                                    >{post.totalRun - post.pass}</TableBodyCell
                                >
                                <TableBodyCell tdClass="idInTable"
                                    >{Math.round(post.passingRate * 10000) /
                                        100}</TableBodyCell
                                >
                            </TableBodyRow>
                        {/each}
                    {/if}
                </TableBody>
            </Table>

            <div style="text-align: center;">
                <br />
                <br />
                <br />

                <div
                    class="chart"
                    bind:clientWidth={width}
                    bind:clientHeight={height}
                >
                    <svg>
                        <!-- y axis -->
                        <g
                            class="axis y-axis"
                            transform="translate(0, {padding.top})"
                        >
                            {#each yTicks as tick}
                                <g
                                    class="tick tick-{tick}"
                                    transform="translate(0, {yScale(tick) -
                                        padding.bottom})"
                                >
                                    <line x2="100%" />
                                    <text y="-4">{tick} %</text>
                                </g>
                            {/each}
                        </g>

                        <!-- x axis -->
                        <g class="axis x-axis">
                            {#each xTicks as tick}
                                <g
                                    class="tick tick-{tick}"
                                    transform="translate({xScale(
                                        tick
                                    )},{height})"
                                >
                                    <line
                                        y1="-{height}"
                                        y2="-{padding.bottom}"
                                        x1="0"
                                        x2="0"
                                    />
                                    {#if tick === 7}<text y="-2"
                                            >{width > 380
                                                ? tick
                                                : formatMobile(tick)} days ago</text
                                        >{:else}<text y="-2"
                                            >{width > 380
                                                ? tick
                                                : formatMobile(tick)}</text
                                        >
                                    {/if}
                                </g>
                            {/each}
                        </g>

                        <!-- data -->
                        <path class="path-area" d={area} />
                        <path class="path-line" d={path} />
                    </svg>
                </div>
            </div>

            <br />
            <br />

            <span class="grid gap-6 mb-6 md:grid-cols-3">
                <div />
                <div>
                    <Label
                        ><span style="font-size: 1rem;">Owner</span>
                        <Select
                            class="mt-2"
                            items={ownerList}
                            bind:value={owner}
                        />
                    </Label>
                    <button
                        on:click={() => {
                            resetOwner();
                            fetchOwner();
                        }}>Search</button
                    >
                </div>
            </span>
        </div>
        <div>
            {#if dataOwner != null}
                {#each dataOwner as testValue (testValue.testName)}
                    <TableBodyRow>
                        <div style="padding-left: 10vw;">
                            <Progressbar
                                progress={`${
                                    Math.round(testValue.passingRate * 10000) /
                                    100
                                }`}
                                color="indigo"
                                labelOutside={testValue.testName}
                                size="h-4"
                            />
                        </div>
                        <br />
                        <br />
                    </TableBodyRow>
                {/each}
            {/if}
        </div>
    </div>
</main>

<style>
    .chart,
    p {
        width: 100%;
        max-width: 70vw;
        margin-left: auto;
        margin-right: auto;
    }

    svg {
        position: relative;
        width: 100%;
        height: 50vh;
        overflow: visible;
    }

    .tick {
        font-size: 0.725em;
        font-weight: 200;
    }

    .tick line {
        stroke: #aaa;
        stroke-dasharray: 2;
    }

    .tick text {
        fill: #666;
        text-anchor: start;
    }

    .tick.tick-0 line {
        stroke-dasharray: 0;
    }

    .x-axis .tick text {
        text-anchor: middle;
    }

    .path-line {
        fill: none;
        stroke: rgb(255, 255, 0);
        stroke-linejoin: round;
        stroke-linecap: round;
        stroke-width: 2;
    }

    .path-area {
        fill: rgba(0, 254, 254, 0.2);
    }
</style>
