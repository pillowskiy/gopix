import Section from "@/components/section";
import styles from "./page.module.scss";
import { BackgroundGradient } from "@/components/backdroung-texture";
import { getDiscoverImages, ImageSort } from "@/shared/actions/images";

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
            <div className={styles.card} key={item.id}>
              <div
                style={{
                  // @ts-ignore
                  "--aspect-ratio":
                    item.properties.width / item.properties.height,
                }}
                className={styles.image}
              >
                <img
                  src={`https://f003.backblazeb2.com/file/s3gopix/${item.path}`}
                  alt={item.path}
                />
              </div>
            </div>
          ))}
        </section>
      </Section>
    </>
  );
}
