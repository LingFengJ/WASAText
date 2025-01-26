<script>
export default {
  name: 'NewGroupView',
  data() {
    return {
      groupName: '',
      members: [],
      newMember: '',
      firstMessage: '',
      error: null,
      loading: false,
      userSuggestions: [],
      showSuggestions: false
    }
  },

  created() {
    document.addEventListener('click', this.handleClickOutside);
  },
  methods: {
    async addMember() {
      if (!this.newMember) return;

      try {
        // Check if user exists
        const response = await this.$axios.get(
          `/users/search?query=${encodeURIComponent(this.newMember)}`,
          {
            headers: {
              'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`
            }
          }
        );

        // Check if we found an exact match
        const userExists = response.data.some(user => 
            user.username === this.newMember
        );

        if (userExists) {
            if (!this.members.includes(this.newMember)) {
                this.members.push(this.newMember);
                this.newMember = '';
                this.error = null;
            }
        } else {
            this.error = 'User does not exist';
            setTimeout(() => this.error = null, 3000); // Clear error after 3s
        }
      } catch (error) {
          console.error('Error checking user:', error);
          this.error = 'Failed to verify user';
      }
    },
    removeMember(member) {
      this.members = this.members.filter(m => m !== member)
    },
    async createGroup() {
      if (!this.groupName || !this.firstMessage || this.members.length === 0) {
        this.error = 'Please fill in all fields and add at least one member'
        return
      }

      this.loading = true
      try {
        const response = await this.$axios.post('/conversations//messages',           
        {
          groupName: this.groupName,
          members: this.members,
          content: this.firstMessage,
          type: 'text'
        },
        {
          headers: {
            'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`,
            'Content-Type': 'application/json'
          },
        })

        const data = response.data
        this.$router.push(`/conversations/${data.conversationId}`)

      } catch (error) {
        this.error = 'Failed to create group'
        console.error('Error:', error)
      } finally {
        this.loading = false
      }
    },

    async searchUser(value) {
      try {
          const response = await this.$axios.get(
              `/users/search?query=${encodeURIComponent(value)}`,
              {
                  headers: {
                      'Authorization': `Bearer ${sessionStorage.getItem('authToken')}`
                  }
              }
          );
          // Filter out current user and already added members
          this.userSuggestions = response.data.filter(user => 
              user.username !== sessionStorage.getItem('username') && 
              !this.members.includes(user.username)
          );
          this.showSuggestions = true;
      } catch (error) {
          console.error('Error searching users:', error);
      }
    },

    //we want that addMember to only accept valid users
    selectUser(username) {
      this.members.push(username);
      this.newMember = '';
      this.showSuggestions = false;
      this.userSuggestions = [];
    },

    handleClickOutside(event) {
        if (!event.target.closest('.input-group') && !event.target.closest('.user-suggestions')) {
            this.showSuggestions = false;
        }
    },
  }
}
</script>

<template>
  <div class="new-group-view p-4">
    <h2>Create New Group</h2>
    
    <div v-if="error" class="alert alert-danger">
      {{ error }}
    </div>

    <form @submit.prevent="createGroup">
      <div class="mb-3">
        <label class="form-label">Group Name:</label>
        <input 
          type="text" 
          class="form-control" 
          v-model="groupName" 
          required
          placeholder="Enter group name"
        >
      </div>

      <div class="mb-3">
        <label class="form-label">Add Members:</label>
        <div class="input-group">
          <input 
            type="text" 
            class="form-control" 
            v-model="newMember"
            @input="searchUser(newMember)"
            @focus="searchUser(newMember)"
            placeholder="Enter username"
          >
          <button 
            type="button" 
            class="btn btn-secondary" 
            @click="addMember"
          >
            Add
          </button>
        </div>

        <!-- User suggestions dropdown -->
        <div v-if="showSuggestions && userSuggestions.length > 0" class="user-suggestions">
            <div 
                v-for="user in userSuggestions" 
                :key="user.id"
                class="suggestion-item"
                @click="selectUser(user.username)"
            >
                <i class="bi bi-person-circle me-2"></i>
                {{ user.username }}
            </div>
        </div>
      </div>

      <div class="mb-3" v-if="members.length > 0">
        <label class="form-label">Members:</label>
        <div class="list-group">
          <div 
            v-for="member in members" 
            :key="member" 
            class="list-group-item d-flex justify-content-between align-items-center"
          >
            {{ member }}
            <button 
              type="button" 
              class="btn btn-sm btn-danger" 
              @click="removeMember(member)"
            >
              Remove
            </button>
          </div>
        </div>
      </div>

      <div class="mb-3">
        <label class="form-label">First Message:</label>
        <textarea 
          class="form-control" 
          v-model="firstMessage" 
          required
          rows="3"
          placeholder="Type your first message..."
        ></textarea>
      </div>

      <button 
        type="submit" 
        class="btn btn-primary" 
        :disabled="loading"
      >
        {{ loading ? 'Creating...' : 'Create Group' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.new-group-view {
  max-width: 800px;
  margin: 0 auto;
}

.list-group {
  margin-top: 10px;
}

.list-group-item {
  margin-bottom: 5px;
}

.user-suggestions {
    position: absolute;
    width: 100%;
    max-height: 200px;
    overflow-y: auto;
    background: white;
    border: 1px solid #ddd;
    border-radius: 4px;
    z-index: 1000;
    margin-top: 2px;
}

.suggestion-item {
    padding: 8px 16px;
    cursor: pointer;
}

.suggestion-item:hover {
    background-color: #f8f9fa;
}
</style>