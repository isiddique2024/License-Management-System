import { useRef } from "react";
import { effect, signal } from "@preact/signals-react";

const showContent = signal<boolean>(false);
const isFadingIn = signal<boolean>(false);

const useLoadingEffect = (data: any, isLoading: boolean, isLoadingDone: any) => {
  const effectExecuted = useRef(false);

  effect(() => {
    if (!effectExecuted.current && data && !isLoading) {
      console.log("license page loading effect");
      isFadingIn.value = true;
      setTimeout(() => {
        showContent.value = true;
        setTimeout(() => {
          isLoadingDone.value = true;
        }, 500);
      }, 1000);

      effectExecuted.current = true;
    }
  });
};

export default useLoadingEffect;
export { showContent, isFadingIn };
