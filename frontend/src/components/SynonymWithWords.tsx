import React from 'react';

interface WordWithSynonymsProps {
  synonym: string;
  words: string[];
}

const SynonymWithWords: React.FC<WordWithSynonymsProps> = ({
  words,
  synonym,
}) => {
  return (
    <div className="text-center flex gap-4 flex-col">
      <h1 className="text-4xl">
        Synonym: <span className="text-primary">{synonym}</span>
      </h1>
      <div className="flex justify-center gap-2">
        <p>{'['}</p>
        {words.map((word, index) => (
          <React.Fragment key={index}>
            {index > 0 && ', '}
            <span key={index} className="text-md text-secondary">
              {word}
            </span>
          </React.Fragment>
        ))}
        <p>{']'}</p>
      </div>
    </div>
  );
};

export default SynonymWithWords;
