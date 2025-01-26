<script>
export default {
    data() {
        return {
            username: null,
            authToken: null,
            conversationsPath: "/conversations"
        }
    },
    watch: {
        '$route': {
            immediate: true,
            handler(to) {
                this.username = sessionStorage.getItem('username');
                this.authToken = sessionStorage.getItem('authToken');
                if (this.authToken === null && to.path !== '/login') {
                    this.$router.push('/login');
                }
            }
        }
    },
    methods: {
        logout() {
            localStorage.clear();
            sessionStorage.clear();
            this.$router.push('/login');
            location.reload();
        },
        refresh() {
            location.reload();
        },
        createNewChat() {
            this.$router.push('/new-conversation');
        },
        createNewGroup() {
            this.$router.push('/new-group');
        }
    }
}
</script>

<template>
    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
        <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAText</a>
        <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse"
            data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>
    </header>

    <div class="container-fluid">
        <div class="row">
            <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
                <div v-if="authToken !== null" class="position-sticky pt-3 sidebar-sticky">
                    <h6
                        class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                        <span>Conversations</span>
                    </h6>
                    <ul class="nav flex-column">
                        <li class="nav-item">
                            <RouterLink :to="conversationsPath" class="nav-link">
                                <i class="bi bi-chat-dots"></i>
                                All Conversations
                            </RouterLink>
                        </li>
                        <li class="nav-item">
                            <a href="#" class="nav-link" @click="createNewChat">
                                <i class="bi bi-plus-circle"></i>
                                New Chat
                            </a>
                        </li>
                        <li class="nav-item">
                            <a href="#" class="nav-link" @click="createNewGroup">
                                <i class="bi bi-people"></i>
                                New Group
                            </a>
                        </li>
                    </ul>

                    <h6
                        class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                        <span>Settings</span>
                    </h6>
                    <ul class="nav flex-column">
                        <li class="nav-item">
                            <RouterLink to="/profile" class="nav-link">
                                <i class="bi bi-person"></i>
                                Profile: {{ username }}
                            </RouterLink>
                        </li>
                        <li class="nav-item">
                            <a href="#" class="nav-link" @click="logout">
                                <i class="bi bi-box-arrow-right"></i>
                                Logout
                            </a>
                        </li>
                    </ul>
                </div>
            </nav>

            <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                <RouterView />
            </main>
        </div>
    </div>
</template>

<style scoped>
.nav-link i {
    margin-right: 8px;
}
</style>