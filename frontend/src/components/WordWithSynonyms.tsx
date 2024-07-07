import React from 'react';

interface WordWithSynonymsProps {
  word: string;
  synonyms: string[];
}

const WordWithSynonyms: React.FC<WordWithSynonymsProps> = ({
  word,
  synonyms,
}) => {
  return (
    <div className="text-center flex gap-4 flex-col">
      <h1 className="text-4xl">
        Word: <span className="text-primary">{word}</span>
      </h1>
      <div className="flex justify-center gap-2 flex-wrap">
        <p>{'['}</p>
        {synonyms.map((synonym, index) => (
          <React.Fragment key={index}>
            {index > 0 && ','}
            <span key={index} className="text-md text-green-300">
              {synonym}
            </span>
          </React.Fragment>
        ))}
        <p>{']'}</p>
      </div>
    </div>
  );
};

export default WordWithSynonyms;
