<!-- Upload CSV -->
<script>
import { Button } from 'flowbite-svelte';
import { P } from 'flowbite-svelte'
</script>

<br>
<br>
<br>

<div class="description">
    <P class="mb-3 md:text-xl" weight="light" size="lg" color="text-gray-500 dark:text-gray-400" firstupper
    >This page is a temporary solution for submitting the data from google's lab. Once we have the access of their DB, we might be able to sql the results more efficiently.
    </P>
</div>
<br>
<br>
<br>

<div style="text-align: center; color:aquamarine; font-size: 1rem; position:relative; padding-left: 30vw; width: 70vw">
    <form
        action="http://10.240.102.16:8082/uploadCSV"
        method="POST"
        enctype="multipart/form-data"
        style="display: flex;"
    >
        upload file:<input name="stainlessData" type="file" /><br />
        <Button shadow="blue" gradient color="blue" type="submit">submit</Button>
    </form>
</div>

<style>
    .description{
        padding-left: 30vw;
        width:70vw;
    }
</style>
