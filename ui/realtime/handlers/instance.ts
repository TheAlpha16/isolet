export const instanceHandlers = {
  "instance.created": (payload: any) => {
    console.log("Instance created:", payload);
    // TODO: implement state updates, e.g. mutate react-query cache or redux store
  },
  "instance.updated": (payload: any) => {
    console.log("Instance updated:", payload);
  },
  "instance.deleted": (payload: any) => {
    console.log("Instance deleted:", payload);
  },
};
