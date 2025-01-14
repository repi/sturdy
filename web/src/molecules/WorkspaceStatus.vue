<template>
  <Tooltip v-if="workspace.statuses.length > 0" :disabled="!isStale">
    <template #tooltip> Draft changed since the last run </template>
    <template #default>
      <StatusDetails
        :statuses="workspace.statuses"
        :class="{
          'opacity-50': isStale,
        }"
      />
    </template>
  </Tooltip>

  <Button
    v-if="(workspace.statuses.length == 0 && ciEnabled) || isStale"
    :icon="terminalIcon"
    :spinner="triggering"
    class="border-0 -ml-3"
    @click="onTriggerClicked"
    >Trigger CI</Button
  >
</template>

<script lang="ts">
import { defineComponent, type PropType, toRefs } from 'vue'
import { gql } from '@urql/vue'

import StatusDetails, { STATUS_FRAGMENT } from '../components/statuses/StatusDetails.vue'
import Button from '../atoms/Button.vue'
import Tooltip from '../atoms/Tooltip.vue'
import { TerminalIcon } from '@heroicons/vue/outline'

import type { WorkspaceStatus_WorkspaceFragment } from './__generated__/WorkspaceStatus'
import { IntegrationProvider } from '../__generated__/types'
import { useUpdatedWorkspacesStatuses } from '../subscriptions/useUpdatedWorkspacesStatuses'
import { useTriggerInstantIntegration } from '../mutations/useTriggerInstantIntegration'

export const WORKSPACE_FRAGMENT = gql`
  fragment WorkspaceStatus_Workspace on Workspace {
    id
    statuses {
      id
      stale
      ...Status
    }
    codebase {
      id
      integrations {
        id
        provider
      }
    }
  }
  ${STATUS_FRAGMENT}
`

const ciProviders = [IntegrationProvider.Buildkite]

export default defineComponent({
  components: { StatusDetails, Button, Tooltip },
  props: {
    workspace: { type: Object as PropType<WorkspaceStatus_WorkspaceFragment>, required: true },
  },
  setup(props) {
    const { workspace } = toRefs(props)
    const ciEnabled = workspace.value.codebase.integrations.some(({ provider }) =>
      ciProviders.includes(provider)
    )
    if (ciEnabled) {
      useUpdatedWorkspacesStatuses([workspace.value.id])
    }
    const triggerInstantIntegration = useTriggerInstantIntegration()
    return { triggerInstantIntegration, terminalIcon: TerminalIcon, ciEnabled }
  },
  data() {
    return {
      triggering: false,
    }
  },
  computed: {
    isStale() {
      return this.workspace.statuses.some(({ stale }) => stale)
    },
  },
  methods: {
    onTriggerClicked() {
      this.triggering = true
      this.triggerInstantIntegration({
        workspaceID: this.workspace.id,
      }).finally(() => (this.triggering = false))
    },
  },
})
</script>
