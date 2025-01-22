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
      loading: false
    }
  },
  methods: {
    addMember() {
      if (this.newMember && !this.members.includes(this.newMember)) {
        this.members.push(this.newMember)
        this.newMember = ''
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
          type: 'group',
          groupName: this.groupName,
          members: this.members,
          content: this.firstMessage,
          messageType: 'text'
        },
        {
          method: 'POST',
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
    }
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
</style>