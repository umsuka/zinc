<template>
  <div class="column index-menu">
    <div>
      <q-select
        v-model="selectedIndex"
        :options="options"
        filled
        dense
        use-input
        input-debounce="0"
        label="Select Index First"
        behavior="menu"
        class="q-mt-md q-mb-sm"
        @filter="filterFn"
        @update:model-value="selectFn"
      >
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey"> No results </q-item-section>
          </q-item>
        </template>
      </q-select>
    </div>
    <div class="index-table">
      <q-table
        v-model:selected="selectedFields"
        :rows="indexFields"
        row-key="name"
        :filter="filterField"
        :filter-method="filterFieldFn"
        :pagination="{ rowsPerPage: 999 }"
        dense
        hide-header
        hide-bottom
        selection="multiple"
        class="field-table"
        @row-click="clickFieldFn"
        @update:selected="selectedFieldFn"
      >
        <template #top-right>
          <q-input
            v-model="filterField"
            filled
            borderless
            dense
            clearable
            debounce="1"
            placeholder="Search for a field"
          >
            <template #append>
              <q-icon name="search" />
            </template>
          </q-input>
        </template>
      </q-table>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref } from "vue";
import indexService from "../../services/index";

export default defineComponent({
  name: "ComponentSearchIndexSelect",
  props: {
    data: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ["updated"],
  setup(props, { emit }) {
    const getIndexData = (field) => props.data[field];
    const selectedIndex = ref(getIndexData("name"));
    const selectedFields = ref(getIndexData("columns"));
    const indexList = ref([]);
    const indexFields = ref([]);
    const cachedFields = ref({});
    const options = ref([]);

    const defaultFields = () => {
      return [{ name: "_id" }, { name: "_index" }, { name: "_score" }];
    };

    // index operation
    const filterFn = (val, update) => {
      if (val === "") {
        update(() => {
          options.value = indexList.value;
        });
        return;
      }

      update(() => {
        const needle = val.toLowerCase();
        options.value = indexList.value.filter((v) =>
          v.value.toLowerCase().includes(needle)
        );
      });
    };

    const selectFn = (index) => {
      if (!index || !index.value) {
        return;
      }
      selectedFields.value = [];
      indexFields.value = defaultFields();
      cachedFields.value = {};

      emit("updated", {
        name: index.value,
        columns: [],
      });
    };

    const setIndexFields = (fields) => {
      for (var k in fields) {
        if (cachedFields.value[fields[k]]) {
          continue;
        }
        indexFields.value.push({ name: fields[k] });
        cachedFields.value[fields[k]] = true;
      }
    };

    // fields operation
    const filterField = ref("");
    const filterFieldFn = (rows, terms) => {
      var filtered = [];
      terms = terms.toLowerCase();
      for (var i = 0; i < rows.length; i++) {
        if (rows[i]["name"].toLowerCase().includes(terms)) {
          filtered.push(rows[i]);
        }
      }
      return filtered;
    };
    const clickFieldFn = (evt, row, index) => {
      if (selectedFields.value.includes(row)) {
        selectedFields.value = selectedFields.value.filter(
          (v) => v.name !== row.name
        );
      } else {
        selectedFields.value.push(row);
      }
      emit("updated", {
        name: selectedIndex.value.value,
        columns: selectedFields.value.map((v) => v.name),
      });
    };
    const selectedFieldFn = () => {
      emit("updated", {
        name: selectedIndex.value.value,
        columns: selectedFields.value.map((v) => v.name),
      });
    };

    // get the list of indices from server when the component is mounted
    const getIndexList = () => {
      indexService.list().then((res) => {
        res.data.map((item) => {
          indexList.value.push({
            label: item.name,
            value: item.name,
          });
        });
      });
    };

    getIndexList();

    return {
      selectedIndex,
      selectedFields,
      options,
      filterFn,
      selectFn,
      indexFields,
      cachedFields,
      filterField,
      filterFieldFn,
      clickFieldFn,
      selectedFieldFn,
      setIndexFields,
    };
  },
});
</script>

<style lang="scss">
.index-menu {
  width: 220px;
  .index-table {
    width: 100%;
  }
  .field-table {
    width: 100%;
  }
}
</style>
