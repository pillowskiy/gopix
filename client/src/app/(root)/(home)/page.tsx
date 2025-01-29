import Section from "@/components/section";
import styles from "./page.module.scss";
import { BackgroundGradient } from "@/components/backdroung-texture";
import { getDiscoverImages, ImageSort } from "@/shared/actions/images";
import { ImageCard } from "@/components/image-card";

export default async function Home() {
  const discover = await getDiscoverImages(2, 50, ImageSort.Newest);

  return (
    <>
      <div className={styles.header}>
        <h1 className={styles.headerTitle}>
          Gopix UI brings headless UI components to life, inspired by
          HeadlessUI, seamlessly blending innovation with flexibility.
        </h1>
      </div>
      <BackgroundGradient />
      <Section className={styles.page} container={false}>
        <section className={styles.container}>
          {discover.items.map((item) => (
            <ImageCard key={item.id} image={item} />
          ))}
        </section>
      </Section>
    </>
  );
}
