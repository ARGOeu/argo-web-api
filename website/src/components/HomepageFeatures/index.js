import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
    {
      title: 'Define Tenants, topologies and feeds',
      Svg: require('@site/static/img/undraw_tenants.svg').default,
      description: (
        <>
          Define tenants to be monitored along with their topologies and feeds
        </>
      ),
    },
    {
      title: 'Create reports and profiles',
      Svg: require('@site/static/img/undraw_documents.svg').default,
      description: (
        <>
          For each tenant create reports with metric and aggregation profiles to kickoff computations and provide results
        </>
      ),
    },
    {
      title: 'Explore A/R and status results',
      Svg: require('@site/static/img/undraw_charts.svg').default,
      description: (
        <>
          Exlpore A/R and status results per groups, endpoints and metrics
        </>
      ),
    }
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
