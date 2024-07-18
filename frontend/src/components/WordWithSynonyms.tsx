import React from 'react'

interface WordWithSynonymsProps {
  word: string
  synonyms: string[]
}

const WordWithSynonyms: React.FC<WordWithSynonymsProps> = ({
  word,
  synonyms,
}) => {
  return (
    <div className="text-center flex gap-4 flex-col">
      <h1 className="text-4xl text-primary">{word}</h1>
      <div className="flex justify-center flex-wrap">
        <p className="text-white">{'['}&nbsp;</p>
        {synonyms.map((synonym, index) => (
          <span key={index} className="text-md text-green-300">
            {index > 0 && <span className="text-white">,&nbsp;</span>}
            {synonym}
          </span>
        ))}
        <p className="text-white">&nbsp;{']'}</p>
      </div>
    </div>
  )
}

export default WordWithSynonyms
