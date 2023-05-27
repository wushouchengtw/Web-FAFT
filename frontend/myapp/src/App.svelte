<script>
  import {
    Navbar,
    NavBrand,
    NavLi,
    NavUl,
    NavHamburger,
    Dropdown,
    DropdownItem,
    Chevron,
    DropdownDivider,
  } from "flowbite-svelte";
  import { Router, Route } from "svelte-navigator";
  import LocalTest from "./localTest.svelte";
  import UploadCsv from "./uploadCSV.svelte";
  import Home from "./home.svelte";
  import RedMine from "./redMine.svelte";
  import Knowledge from "./knowledge.svelte";
  import Stainless from "./stainless.svelte";
  import CaseStatus from "./caseStatus.svelte";
</script>

<Router>
  <Navbar let:hidden let:toggle navClass="full-width" color="none">
    <NavBrand href="/">
      <span
        class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
        style="font-size: 3rem;"
        >Testthaus
      </span>
    </NavBrand>
    <NavHamburger on:click={toggle} />
    <NavUl {hidden}>
      <NavLi href="/home" active={true} activeClass="navFont2">Home</NavLi>
      <NavLi
        id="nav-menu1"
        class="cursor-pointer"
        active={true}
        activeClass="navFont2"><Chevron aligned>Stainless</Chevron></NavLi
      >
      <NavLi href="localTest" active={true} activeClass="navFont2"
        >localTest</NavLi
      >
      <NavLi href="/redMine" active={true} activeClass="navFont2">RedMine</NavLi
      >
      <NavLi href="/knowledge" active={true} activeClass="navFont2"
        >Knowledge</NavLi
      >
      <Dropdown triggeredBy="#nav-menu1" class="w-44 z-20">
        <DropdownItem href="/caseStatus">Dashboard</DropdownItem>
        <DropdownItem href="/stainless">Search test</DropdownItem>
        <DropdownDivider />
        <DropdownItem href="/uploadCSV">Upload CSV</DropdownItem>
      </Dropdown>
    </NavUl>
  </Navbar>

  <Route path="home">
    <Home />
  </Route>
  <Route path="caseStatus">
    <CaseStatus />
  </Route>
  <Route path="stainless">
    <Stainless />
  </Route>
  <Route path="uploadCSV">
    <UploadCsv />
  </Route>
  <Route path="localTest">
    <LocalTest />
  </Route>
  <Route path="redMine">
    <RedMine />
  </Route>
  <Route path="knowledge">
    <Knowledge />
  </Route>
</Router>
